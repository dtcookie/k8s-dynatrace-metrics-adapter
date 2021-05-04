package metrics

import (
	"encoding/json"
)

type QueryResult struct {
	TimeseriesID string      `json:"timeseriesId"`
	DataResult   *DataResult `json:"dataResult"`
	TSQueried    int64       `json:"-"`
	TSRequested  int64       `json:"-"`
}

func (tdpqr *QueryResult) GetID() string {
	return tdpqr.TimeseriesID
}

func (tdpqr *QueryResult) GetTSQueried() int64 {
	return tdpqr.TSQueried
}

func (tdpqr *QueryResult) SetTSQueried(v int64) {
	tdpqr.TSQueried = v
}

func (tdpqr *QueryResult) GetTSRequested() int64 {
	return tdpqr.TSRequested
}

func (tdpqr *QueryResult) SetTSRequested(v int64) {
	tdpqr.TSRequested = v
}

type DataResult struct {
	TimeseriesID    string                `json:"timeseriesId"`
	DataPoints      map[string]DataPoints `json:"dataPoints"`
	AggregationType *AggregationType      `json:"aggregationType"`
	Entities        map[string]string     `json:"entities"`
}

type DataPoints []*DataPoint

func (dp DataPoints) Reduce() *DataPoint {
	if len(dp) == 0 {
		return nil
	}
	var result *DataPoint
	for _, entry := range dp {
		if entry == nil {
			continue
		}
		if entry.Value == nil {
			continue
		}
		result = entry
	}
	return result
}

type DataPoint struct {
	TimeStamp   int64
	Value       *float64
	StringValue *string
}

func (dp *DataPoint) UnmarshalJSON(data []byte) error {
	rawData := []json.RawMessage{}
	if err := json.Unmarshal(data, &rawData); err != nil {
		return err
	}
	if err := json.Unmarshal(rawData[0], &dp.TimeStamp); err != nil {
		return err
	}
	if err := json.Unmarshal(rawData[1], &dp.Value); err != nil {
		if err := json.Unmarshal(rawData[1], &dp.StringValue); err != nil {
			return err
		}
	}
	return nil
}
