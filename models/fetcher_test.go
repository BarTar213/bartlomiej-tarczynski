package models

import (
	"reflect"
	"testing"
)

const (
	validUrl = "https://httpbin.org/range/15"
	invalidUrl = "abcUrl"
)


func TestFetcher_Reset(t *testing.T) {
	type fields struct {
		Id       int
		Url      string
		Interval int
		JobId    int
	}
	tests := []struct {
		name   string
		fields fields
	}{
		{
			name: "positive_reset",
			fields: fields{
				Id:       12,
				Url:      validUrl,
				Interval: 15,
				JobId:    18,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				Id:       tt.fields.Id,
				Url:      tt.fields.Url,
				Interval: tt.fields.Interval,
				JobId:    tt.fields.JobId,
			}
			f.Reset()
			if !reflect.DeepEqual(f, &Fetcher{}) {
				t.Errorf("Reset() = %v, wantErr %v", f, &Fetcher{})
			}
		})
	}
}

func TestFetcher_Validate(t *testing.T) {
	type fields struct {
		Id       int
		Url      string
		Interval int
		JobId    int
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name:    "positive_validate",
			fields:  fields{
				Url:      validUrl,
				Interval: 5,
			},
			wantErr: false,
		},
		{
			name:    "negative_validate_empty_url_error",
			fields:  fields{
				Url:      "",
				Interval: 5,
			},
			wantErr: true,
		},
		{
			name:    "negative_validate_invalid_url_error",
			fields:  fields{
				Url:      invalidUrl,
				Interval: 5,
			},
			wantErr: true,
		},
		{
			name:    "negative_validate_invalid_interval_error",
			fields:  fields{
				Url:      validUrl,
				Interval: -5,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fetcher{
				Id:       tt.fields.Id,
				Url:      tt.fields.Url,
				Interval: tt.fields.Interval,
				JobId:    tt.fields.JobId,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
