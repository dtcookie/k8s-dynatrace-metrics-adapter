package metrics

type AggregationTypes []AggregationType

func (at AggregationTypes) Contains(v AggregationType) bool {
	for _, elem := range at {
		if elem == v {
			return true
		}
	}
	return false
}

type AggregationType string

var AllAggregationTypes = struct {
	Avg        AggregationType
	Count      AggregationType
	Max        AggregationType
	Median     AggregationType
	Min        AggregationType
	Percentile AggregationType
	Sum        AggregationType
}{
	Avg:        AggregationType("AVG"),
	Count:      AggregationType("COUNT"),
	Max:        AggregationType("MAX"),
	Median:     AggregationType("MEDIAN"),
	Min:        AggregationType("MIN"),
	Percentile: AggregationType("PERCENTILE"),
	Sum:        AggregationType("SUM"),
}
