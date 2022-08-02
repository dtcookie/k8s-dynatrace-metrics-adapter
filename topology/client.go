package topology

import (
	"encoding/json"
	"fmt"
	"net/url"
	"strings"

	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/cache"
	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/metrics"
	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/rest"
)

type Client struct {
	restClient rest.Client
}

func NewClient(config *rest.Config, baseURL string, token string) *Client {
	return &Client{restClient: rest.NewClient(config, baseURL, rest.NewCredentials(token))}
}

func (tc *Client) GetMetrics() ([]cache.Item, error) {
	var err error
	var data []byte
	if data, err = tc.restClient.GET("/api/v1/timeseries", 200); err != nil {
		return nil, err
	}
	metrics := []*metrics.Metric{}
	if err := json.Unmarshal(data, &metrics); err != nil {
		return nil, err
	}
	result := make([]cache.Item, len(metrics))
	for _, item := range metrics {
		result = append(result, item)
	}
	return result, nil
}

func isAPIv1(id string) bool {
	return strings.HasPrefix(id, "v1:") || strings.HasPrefix(id, "com.dynatrace.builtin:") || strings.HasPrefix(id, "custom")
}

func (tc *Client) getMetricAPIv1(id string) (cache.Item, error) {
	var err error
	var data []byte

	if data, err = tc.restClient.GET(fmt.Sprintf("/api/v1/timeseries/%s?includeData=false", id), 200); err != nil {
		return nil, err
	}
	var metric metrics.Metric
	if err := json.Unmarshal(data, &metric); err != nil {
		return nil, err
	}

	return &metric, nil
}

func (tc *Client) getMetricAPIv2(metricId string) (cache.Item, error) {
	var err error
	var data []byte

	if data, err = tc.restClient.GET(fmt.Sprintf("/api/v2/metrics/%s", metricId), 200); err != nil {
		return nil, err
	}

	var metric metrics.Metric
	var metricV2 metrics.MetricV2

	if err := json.Unmarshal(data, &metricV2); err != nil {
		return nil, err
	}

	metric.ID = metricV2.ID
	metric.Dimensions = []string{}
	for _, dim := range metricV2.Dimensions {
		metric.Dimensions = append(metric.Dimensions, dim.Key)
	}
	metric.AggregationTypes = metricV2.AggregationTypes

	return &metric, nil
}

func (tc *Client) GetMetric(id string, options map[string]interface{}) (cache.Item, error) {
	if strings.Contains(id, ":#") {
		return nil, nil
	}

	if isAPIv1(id) {
		return tc.getMetricAPIv1(id)
	}

	return tc.getMetricAPIv2(id)
}

type MetricList struct {
	cache.ItemBase
	Items []*metrics.Metric
}

func (tc *Client) GetAllMetrics(fakeID string, options map[string]interface{}) (cache.Item, error) {
	result := []*metrics.Metric{}
	var items []cache.Item
	var err error
	if items, err = tc.GetMetrics(); err != nil {
		return nil, err
	}
	for _, item := range items {
		if item == nil {
			continue
		}
		result = append(result, item.(*metrics.Metric))
	}
	return &MetricList{Items: result}, nil
}

func Hide(v interface{}) {}

func (tc *Client) getDataPointsAPIv1(
	timeseriesID string,
	aggregationType *metrics.AggregationType,
	tags map[string]*string,
) (cache.Item, error) {

	var err error
	var data []byte
	var url string
	if aggregationType != nil {
		url = fmt.Sprintf("/api/v1/timeseries/%s?includeData=true&relativeTime=10mins&aggregationType=%v&includeParentIds=true", timeseriesID, *aggregationType)
	} else {
		url = fmt.Sprintf("/api/v1/timeseries/%s?includeData=true&relativeTime=10mins&includeParentIds=true", timeseriesID)
	}
	if len(tags) > 0 {
		for k, v := range tags {
			if v == nil {
				url = fmt.Sprintf("%v&tag=%v", url, k)
			} else {
				url = fmt.Sprintf("%v&tag=%v:%v", url, k, *v)
			}
		}
	}
	if data, err = tc.restClient.GET(url, 200); err != nil {
		return nil, err
	}
	result := metrics.QueryResult{}
	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}
	return &result, nil
}

func (tc *Client) getDataPointsAPIv2(metricID string, aggregationType *metrics.AggregationType) (cache.Item, error) {

	var err error
	var data []byte
	var path = fmt.Sprintf("/api/v2/metrics/query?metricSelector=%s", url.QueryEscape(metricID))

	// time range and resolution
	path = fmt.Sprintf("%v&resolution=m&from=now-5m&to=now", path)

	if data, err = tc.restClient.GET(path, 200); err != nil {
		return nil, err
	}
	result := metrics.QueryResult{}
	resultV2 := metrics.QueryResultV2{}

	if err := json.Unmarshal(data, &resultV2); err != nil {
		return nil, err
	}

	if len(resultV2.Result) == 0 || len(resultV2.Result[0].Data) == 0 {
		return nil, nil
	}

	result.TimeseriesID = resultV2.Result[0].ID
	result.DataResult = &metrics.DataResult{}
	result.DataResult.TimeseriesID = resultV2.Result[0].ID
	result.DataResult.AggregationType = aggregationType
	result.DataResult.Entities = resultV2.Result[0].Data[0].DimensionsMap
	result.DataResult.DataPoints = map[string]metrics.DataPoints{}

	dataPoints := metrics.DataPoints{}
	for index, timestamp := range resultV2.Result[0].Data[0].Timestamps {
		value := resultV2.Result[0].Data[0].Values[index]

		dataPoint := metrics.DataPoint{
			TimeStamp: timestamp,
			Value:     &value,
		}

		dataPoints = append(dataPoints, &dataPoint)
	}

	for _, value := range resultV2.Result[0].Data[0].DimensionsMap {
		result.DataResult.DataPoints[value] = dataPoints
	}

	return &result, nil
}

func (tc *Client) GetDataPoints(timeseriesID string, options map[string]interface{}) (cache.Item, error) {
	var aggregationType *metrics.AggregationType
	var tags map[string]*string

	if options != nil {
		if v, found := options["agg"]; found {
			agg := v.(metrics.AggregationType)
			aggregationType = &agg
		}
		if v, found := options["tags"]; found {
			tags = v.(map[string]*string)
		}
	}

	if isAPIv1(timeseriesID) {
		return tc.getDataPointsAPIv1(timeseriesID, aggregationType, tags)
	}

	return tc.getDataPointsAPIv2(timeseriesID, aggregationType)
}

func (tc *Client) GetService(id string, options map[string]interface{}) (cache.Item, error) {
	var err error
	var data []byte
	if data, err = tc.restClient.GET(fmt.Sprintf("/api/v1/entity/services/%s", id), 200); err != nil {
		return nil, err
	}
	var service Service
	if err := json.Unmarshal(data, &service); err != nil {
		return nil, err
	}
	return &service, nil
}

func (tc *Client) GetProcess(id string, options map[string]interface{}) (cache.Item, error) {
	var err error
	var data []byte
	if data, err = tc.restClient.GET(fmt.Sprintf("/api/v1/entity/infrastructure/processes/%s", id), 200); err != nil {
		return nil, err
	}
	var process Process
	if err := json.Unmarshal(data, &process); err != nil {
		return nil, err
	}
	return &process, nil
}

func (tc *Client) GetProcessGroup(id string, options map[string]interface{}) (cache.Item, error) {
	var err error
	var data []byte
	if data, err = tc.restClient.GET(fmt.Sprintf("/api/v1/entity/infrastructure/process-groups/%s", id), 200); err != nil {
		return nil, err
	}
	var processGroup ProcessGroup
	if err := json.Unmarshal(data, &processGroup); err != nil {
		return nil, err
	}
	return &processGroup, nil
}

func (tc *Client) GetHost(id string, options map[string]interface{}) (cache.Item, error) {
	var err error
	var data []byte
	if data, err = tc.restClient.GET(fmt.Sprintf("/api/v1/entity/infrastructure/hosts/%s", id), 200); err != nil {
		return nil, err
	}
	var host Host
	if err := json.Unmarshal(data, &host); err != nil {
		return nil, err
	}
	return &host, nil
}
