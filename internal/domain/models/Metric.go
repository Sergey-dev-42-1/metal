package models

type MetricValue int64

func (v MetricValue) ToFloat64() float64 {
	return float64(v)
}

type Metric struct {
	Value MetricValue
	Type  string
	Name  string
}
