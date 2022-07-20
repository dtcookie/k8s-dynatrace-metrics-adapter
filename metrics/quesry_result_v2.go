package metrics

type QueryResultV2 struct {
	Result []QueryResultV2Details `json:"result"`
}

type QueryResultV2Details struct {
	ID   string              `json:"metricId"`
	Data []QueryResultV2Data `json:"data"`
}

type QueryResultV2Data struct {
	DimensionsMap map[string]string `json:"dimensionMap"`
	Timestamps    []int64           `json:"timestamps"`
	Values        []float64         `json:"values"`
}
