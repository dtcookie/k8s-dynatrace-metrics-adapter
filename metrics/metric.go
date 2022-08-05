package metrics

import (
	"github.com/dtcookie/goutils"
)

type Metric struct {
	TSQueried        int64            `json:"-"`
	TSRequested      int64            `json:"-"`
	ID               string           `json:"timeseriesId"`
	Dimensions       goutils.Strings  `json:"dimensions"`
	AggregationTypes AggregationTypes `json:"aggregationTypes"`
	DetailedSource   string           `json:"detailedSource"`
}

func (m *Metric) SuggestAggregationType() *AggregationType {
	aggregationType := &AllAggregationTypes.Avg
	if !m.AggregationTypes.Contains(*aggregationType) {
		if len(m.AggregationTypes) > 0 {
			aggregationType = &m.AggregationTypes[0]
		} else {
			aggregationType = nil
		}
	}
	return aggregationType
}

func (m *Metric) Labels() map[string]string {
	return map[string]string{}
}

func (m *Metric) GetID() string {
	return m.ID
}
func (m *Metric) GetTSQueried() int64 {
	return m.TSQueried
}

func (m *Metric) SetTSQueried(v int64) {
	m.TSQueried = v
}

func (m *Metric) GetTSRequested() int64 {
	return m.TSRequested
}

func (m *Metric) SetTSRequested(v int64) {
	m.TSRequested = v
}

type DimensionV2 struct {
	Key  string `json:"key"`
	Name string `json:"displayName"`
	Type string `json:"type"`
}

type MetricV2 struct {
	TSQueried        int64            `json:"-"`
	TSRequested      int64            `json:"-"`
	ID               string           `json:"metricId"`
	AggregationTypes AggregationTypes `json:"aggregationTypes"`
	Dimensions        []DimensionV2    `json:"dimensionDefinitions"`
}
