package service

import (
	"runtime"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_CreateMetricsMap(t *testing.T) {
	// type newPerson struct {
	//     r Relationship
	//     p Person
	// }
	tests := []struct {
		name    string
		wantErr bool
	}{
		{
			name:    "Creates proper map",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ms := runtime.MemStats{}
			res := createMetricsMap(ms)
			assert.Equal(t, len(res), 8)
		})
	}
}
