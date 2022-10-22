package collector

import (
	"github.com/david-vtuk/prometheus-rancher-exporter/query/rancher"
	"reflect"
	"testing"
)

func TestCollect(t *testing.T) {
	type args struct {
		client rancher.Client
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			Collect(tt.args.client)
		})
	}
}

func Test_new(t *testing.T) {
	tests := []struct {
		name string
		want metrics
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := new(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("new() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_resetGaugeMetrics(t *testing.T) {
	type args struct {
		m metrics
	}
	tests := []struct {
		name string
		args args
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resetGaugeMetrics(tt.args.m)
		})
	}
}
