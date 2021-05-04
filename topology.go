package main

import (
	"strings"

	"github.com/dtcookie/goutils"
	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/cache"
	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/metrics"
	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/rest"
	"github.com/dtcookie/k8s-dynatrace-metrics-adapter/topology"

	"github.com/kubernetes-sigs/custom-metrics-apiserver/pkg/provider"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/metrics/pkg/apis/external_metrics"
)

type Topology struct {
	client        *topology.Client
	services      cache.Repo
	hosts         cache.Repo
	processGroups cache.Repo
	processes     cache.Repo
	metrics       cache.Repo
	dataPoints    cache.Repo
	allmetrics    cache.Repo
}

func NewTopology(config *rest.Config, baseURL string, token string) *Topology {
	client := topology.NewClient(config, baseURL, token)
	topology := &Topology{
		client:        client,
		services:      cache.NewRepo("service", 60*30, client.GetService, service404),
		hosts:         cache.NewRepo("host", 60*30, client.GetHost, host404),
		processGroups: cache.NewRepo("processGroup", 60*30, client.GetProcessGroup, processGroup404),
		processes:     cache.NewRepo("process", 60*30, client.GetProcess, process404),
		metrics:       cache.NewRepo("metric", 60*5, client.GetMetric, metric404),
		dataPoints:    cache.NewRepo("dataPoints", 40, client.GetDataPoints, nil),
		allmetrics:    cache.NewRepo("allmetrics", 60*5, client.GetAllMetrics, nil),
	}
	return topology
}

func metric404(e error) cache.Item {
	if re, ok := e.(*rest.Error); ok && re.Code == 404 {
		return &metrics.Metric{ID: "404"}
	}
	return nil
}

func host404(e error) cache.Item {
	if re, ok := e.(*rest.Error); ok && re.Code == 404 {
		return &topology.Host{EntityBase: topology.EntityBase{ID: "404"}}
	}
	return nil
}

func service404(e error) cache.Item {
	if re, ok := e.(*rest.Error); ok && re.Code == 404 {
		return &topology.Service{EntityBase: topology.EntityBase{ID: "404"}}
	}
	return nil
}

func process404(e error) cache.Item {
	if re, ok := e.(*rest.Error); ok && re.Code == 404 {
		return &topology.Process{EntityBase: topology.EntityBase{ID: "404"}}
	}
	return nil
}

func processGroup404(e error) cache.Item {
	if re, ok := e.(*rest.Error); ok && re.Code == 404 {
		return &topology.ProcessGroup{EntityBase: topology.EntityBase{ID: "404"}}
	}
	return nil
}

type externalMetric struct {
	info   provider.ExternalMetricInfo
	labels map[string]string
	value  external_metrics.ExternalMetricValue
}

func (top *Topology) GetDataPoints(id string, tags map[string]*string) ([]externalMetric, error) {
	externalMetrics := []externalMetric{}
	var err error
	var item cache.Item
	if item, err = top.metrics.Get(id, nil); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	metric := item.(*metrics.Metric)
	if metric.ID == "404" {
		return externalMetrics, nil
	}
	aggregationType := metric.SuggestAggregationType()
	options := map[string]interface{}{}
	if tags != nil {
		options["tags"] = tags
	}
	if aggregationType != nil {
		options["agg"] = *aggregationType
	}

	if item, err = top.dataPoints.Get(id, options); err != nil {
		return nil, err
	}
	if item == nil {
		return nil, nil
	}
	result := item.(*metrics.QueryResult)
	if result.DataResult == nil {
		return nil, nil
	}
	if result.DataResult.DataPoints == nil {
		return nil, nil
	}
	dims := goutils.StringMap{}
	for k, dataPoints := range result.DataResult.DataPoints {
		if len(dataPoints) == 0 {
			continue
		}
		dataPoint := dataPoints.Reduce()
		if dataPoint == nil {
			continue
		}
		iv := int(*dataPoint.Value * 1000.0)
		iv64 := int64(iv)
		quantity := resource.NewMilliQuantity(iv64, resource.DecimalSI)

		if strings.Contains(k, ",") {
			parts := strings.Split(k, ",")
			for _, part := range parts {
				part = strings.TrimSpace(part)
				var key string
				var value string
				if strings.Contains(part, "=") {
					kv := strings.Split(part, "=")
					key = strings.TrimSpace(kv[0])
					value = strings.TrimSpace(kv[1])
				} else {
					key = metric.Dimensions.GetPrefixFor(part)
					if entity, err := top.Get(part); (err == nil) && (entity != nil) {
						dims.Consume(entity.(topology.Entity).Labels())
					}
					value = part
				}
				value = goutils.Unquote(value)
				dims[key] = value
			}
		} else {
			key := metric.Dimensions.GetPrefixFor(k)
			if entity, err := top.Get(k); (err == nil) && (entity != nil) {
				dims.Consume(entity.(topology.Entity).Labels())
			}
			value := goutils.Unquote(k)
			dims[key] = value
		}

		if len(tags) > 0 {
			for k, v := range tags {
				if v == nil {
					dims[k] = ""
				} else {
					dims[k] = *v
				}
			}
		}

		externalMetric := externalMetric{
			info: provider.ExternalMetricInfo{
				Metric: id,
			},
			labels: dims,
			value: external_metrics.ExternalMetricValue{
				MetricName:   id,
				MetricLabels: dims,
				Value:        *quantity,
			},
		}
		externalMetrics = append(externalMetrics, externalMetric)
	}
	return externalMetrics, nil
}

func (top *Topology) GetMetrics() ([]*metrics.Metric, error) {
	var list cache.Item
	var err error
	if list, err = top.allmetrics.Get("----", nil); err != nil {
		return nil, err
	}
	if list == nil {
		return []*metrics.Metric{}, nil
	}
	items := list.(*topology.MetricList).Items
	return items, nil
}

func (top *Topology) Get(id string) (cache.Item, error) {
	if strings.HasPrefix(id, "HOST-") {
		return top.hosts.Get(id, nil)
	} else if strings.HasPrefix(id, "SERVICE-") {
		return top.services.Get(id, nil)
	} else if strings.HasPrefix(id, "PROCESS_GROUP_INSTANCE-") {
		return top.processes.Get(id, nil)
	} else if strings.HasPrefix(id, "PROCESS_GROUP-") {
		return top.processGroups.Get(id, nil)
	}
	return nil, nil
}
