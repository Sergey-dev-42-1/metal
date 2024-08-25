package models

import (
	"fmt"
)

type MetricValue float64

func (v MetricValue) ToInt64() int64 {
	return int64(v)
}
func (v MetricValue) ToString() string {
	return fmt.Sprintf("%d", v.ToInt64())
}
func (v MetricValue) ToStringFloat() string {
	return fmt.Sprintf("%g", v)
}

type Metric struct {
	Value MetricValue
	Type  string
	Name  string
}
type MetricServiceClient struct {
	Value MetricValue
	Type  string
	Name  string
}
