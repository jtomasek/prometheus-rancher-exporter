package rancher

import (
	"reflect"
	"testing"
)

func TestClient_GetBackups(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []backup
		wantErr bool
	}{
		{"test-1", testClient, []backup{{
			Name:            "test-recurring",
			ResourceSetName: "rancher-resource-set",
			RetentionCount:  10,
			BackupType:      "Recurring",
			Message:         "Completed",
			Filename:        "test-recurring-e3acb0dc-c4f1-4482-83db-66f0141722de-2023-07-19T00-00-00Z.tar.gz",
			LastSnapshot:    "2023-07-19T00:00:09Z",
			NextSnapshot:    "2023-07-20T00:00:00Z",
			StorageLocation: "PV",
		}}, false}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetBackups()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetBackups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetBackups() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetNumberOfBackups(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{"test-1", testClient, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetNumberOfBackups()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNumberOfBackups() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNumberOfBackups() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetNumberOfRestores(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    int
		wantErr bool
	}{
		{"test-1", testClient, 1, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetNumberOfRestores()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNumberOfRestores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetNumberOfRestores() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestClient_GetRestores(t *testing.T) {
	tests := []struct {
		name    string
		fields  fields
		want    []restore
		wantErr bool
	}{
		{"test-1", testClient, []restore{{
			Name:                 "restore-jq9bs",
			Filename:             "one-time-test-2-e3acb0dc-c4f1-4482-83db-66f0141722de-2023-07-19T11-16-41Z.tar.gz",
			Prune:                true,
			StorageLocation:      "PV",
			Message:              "Completed",
			ResoreCompletionTime: "2023-07-19T11:22:07Z",
		}}, false}}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			r := Client{
				Client: tt.fields.Client,
				Config: tt.fields.Config,
			}
			got, err := r.GetRestores()
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRestores() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetRestores() got = %v, want %v", got, tt.want)
			}
		})
	}
}
