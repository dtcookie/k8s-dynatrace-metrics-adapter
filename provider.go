package main

import (
	"os"
	"strings"

	"k8s.io/klog/v2"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/labels"
	"k8s.io/metrics/pkg/apis/external_metrics"

	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/rest"
	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/selection"
)

type testProvider struct {
	topology *Topology
	enabled  bool
}

func NewProvider() provider.ExternalMetricsProvider {
	config := &rest.Config{
		Verbose:  false,
		Debug:    false,
		NoProxy:  false,
		Insecure: true,
	}

	enabled := true
	token := os.Getenv("API_TOKEN")
	if strings.TrimSpace(token) == "" {
		klog.Infof("Environment Variable '%s' not defined. Metric Provider disabled.", "API_TOKEN")
		enabled = false
	}
	baseURL := os.Getenv("BASE_URL")
	if strings.TrimSpace(baseURL) == "" {
		klog.Infof("Environment Variable '%s' not defined. Metric Provider disabled.", "BASE_URL")
		enabled = false
	}
	return &testProvider{
		topology: NewTopology(config, baseURL, token),
		enabled:  enabled,
	}
}

// var allMetrics = []provider.ExternalMetricInfo{}

func (p *testProvider) ListAllExternalMetrics() []provider.ExternalMetricInfo {
	return []provider.ExternalMetricInfo{}
	// if !p.enabled {
	// 	return []provider.ExternalMetricInfo{}
	// }
	// var metricList []*metrics.Metric
	// var err error
	// if metricList, err = p.topology.GetMetrics(); err != nil {
	// 	return []provider.ExternalMetricInfo{}
	// }
	// metricInfos := []provider.ExternalMetricInfo{}
	// for _, metric := range metricList {
	// 	if strings.Contains(metric.ID, "#") {
	// 		continue
	// 	}
	// 	if strings.Contains(metric.ID, "%") {
	// 		continue
	// 	}

	// 	metricInfo := provider.ExternalMetricInfo{Metric: metric.ID}
	// 	// var valueList *external_metrics.ExternalMetricValueList
	// 	// if valueList, err = p.GetExternalMetric("default", labels.NewSelector(), metricInfo); err != nil {
	// 	// 	continue
	// 	// }
	// 	// if valueList == nil || len(valueList.Items) == 0 {
	// 	// 	continue
	// 	// }
	// 	metricInfos = append(metricInfos, metricInfo)
	// }
	// return metricInfos
}

func (p *testProvider) handleError(e error) error {
	if e == nil {
		return e
	}
	switch err := e.(type) {
	case *rest.Error:
		if err.Code == 401 {
			klog.Info("Authentication failed. Check environment variable 'API_TOKEN'. Metric Provider is now disabled.")
			p.enabled = false
		}
		return err
	default:
		return err
	}
}

func (p *testProvider) GetExternalMetric(namespace string, metricSelector labels.Selector, info provider.ExternalMetricInfo) (*external_metrics.ExternalMetricValueList, error) {
	if !p.enabled {
		return nil, nil
	}
	tags := map[string]*string{}
	requirements, _ := metricSelector.Requirements()
	for _, requirement := range requirements {
		key := requirement.Key()
		values := requirement.Values()
		switch requirement.Operator() {
		case selection.Equals:
			if len(values) == 1 {
				tags[key] = &values.List()[0]
			} else {
				tags[key] = nil
			}
		case selection.In:
			tags[key] = nil
		case selection.Exists:
			tags[key] = nil
		case selection.DoubleEquals:
			if len(values) == 1 {
				tags[key] = &values.List()[0]
			} else {
				tags[key] = nil
			}
		default:
		}
	}
	var externalMetrics []externalMetric
	var err error
	if externalMetrics, err = p.topology.GetDataPoints(info.Metric, tags); err != nil {
		klog.Infof("Error: %s", err.Error())
		return nil, p.handleError(err)
	}
	matchingMetrics := []external_metrics.ExternalMetricValue{}
	for _, metric := range externalMetrics {
		if metric.info.Metric == info.Metric {
			if metricSelector.Matches(labels.Set(metric.labels)) {
				metricValue := metric.value
				metricValue.Timestamp = metav1.Now()
				matchingMetrics = append(matchingMetrics, metricValue)
			}
		}
	}
	result := &external_metrics.ExternalMetricValueList{
		Items: matchingMetrics,
	}
	// data, err := json.Marshal(&result)
	// if err != nil {
	// 	klog.Info(err.Error())
	// } else {
	// 	klog.Info(string(data))
	// }
	return result, nil
}
