package rancher

import (
	"reflect"
	"testing"
)

func TestClient_GetNumberofProjects(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{"test-1", testClient, 5, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetNumberofProjects()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNumberofProjects() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNumberofProjects() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetProjectAnnotations(t *testing.T) {

	tests := []struct {
		name    string
		fields  fields
		want    projectAnnotation
		wantErr bool
	}{
		{"test-1", testClient, projectAnnotation{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetProjectAnnotations()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectAnnotations() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var found = false
			for _, v := range got {
				if v == tt.want {
					found = true
				}
			}
			if found == false {
				t.Errorf("GetProjectAnnotations() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetProjectLabels(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    projectLabel
		wantErr bool
	}{
		{"test-1", testClient, projectLabel{"p-wm44k", "Default", "local", "cattle.io/creator", "norman"}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetProjectLabels()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectLabels() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var found = false
			for _, v := range got {
				if v == tt.want {
					found = true
				}
			}
			if found == false {
				t.Errorf("GetProjectLabels() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestClient_GetProjectResourceQuota(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []projectResource
		wantErr bool
	}{
		{"test-1", testClient, []projectResource{}, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetProjectResourceQuota()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetProjectResourceQuota() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetProjectResourceQuota() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_clusterIdToName(t *testing.T) {
	type args struct {
		id string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{"test-1", testClient, args{""}, "", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.clusterIdToName(tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("clusterIdToName() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("clusterIdToName() got = %v, want %v", got, tt.want)
			}
		})
	}
}
