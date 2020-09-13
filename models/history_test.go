package models

import (
	"reflect"
	"testing"
)

func TestHistory_Reset(t *testing.T) {
	type fields struct {
		FetcherId int
		Response  *string
		Duration  float64
		CreatedAt int64
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name:   "positive_reset",
			fields: fields{
				Response: pointer(validUrl),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			h := &History{
				FetcherId: tt.fields.FetcherId,
				Response:  tt.fields.Response,
				Duration:  tt.fields.Duration,
				CreatedAt: tt.fields.CreatedAt,
			}
			h.Reset()
			if !reflect.DeepEqual(h, &History{}) {
				t.Errorf("Reset() = %v, wantErr %v", h, &History{})
			}
		})
	}
}

func pointer(s string) *string{
	return &s
}
