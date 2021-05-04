package topology

import (
	"encoding/json"
	"fmt"
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

func (tc *Client) GetMetric(id string, options map[string]interface{}) (cache.Item, error) {
	if strings.Contains(id, ":#") {
		return nil, nil
	}
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
