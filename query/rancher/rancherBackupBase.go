package rancher

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	backupsGVR  = schema.GroupVersionResource{Group: "resources.cattle.io", Version: "v1", Resource: "backups"}
	restoresGVR = schema.GroupVersionResource{Group: "resources.cattle.io", Version: "v1", Resource: "restores"}
)

type backup struct {
	Name            string
	ResourceSetName string
	RetentionCount  int64
	BackupType      string
	Message         string
	Filename        string
	LastSnapshot    string
	NextSnapshot    string
	StorageLocation string
}

type restore struct {
	Name                 string
	Filename             string
	Prune                bool
	StorageLocation      string
	Message              string
	ResoreCompletionTime string
}

func (r Client) GetNumberOfBackups() (int, error) {

	res, err := r.Client.Resource(backupsGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetNumberOfRestores() (int, error) {

	res, err := r.Client.Resource(restoresGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetBackups() ([]backup, error) {

	var backups []backup
	var backupMessage string
	var backupNextSnapshot string
	var backupLastSnapshot string

	res, err := r.Client.Resource(backupsGVR).List(context.Background(), v1.ListOptions{})

	if err != nil {
		return nil, err
	}

	for _, backupJob := range res.Items {

		backupName, _, err := unstructured.NestedString(backupJob.Object, "metadata", "name")
		if err != nil {
			return nil, err
		}

		backupResourceSetName, _, err := unstructured.NestedString(backupJob.Object, "spec", "resourceSetName")
		if err != nil {
			return nil, err
		}

		backupRetentionCount, _, err := unstructured.NestedInt64(backupJob.Object, "spec", "retentionCount")
		if err != nil {
			return nil, err
		}

		backupType, _, err := unstructured.NestedString(backupJob.Object, "status", "backupType")
		if err != nil {
			return nil, err
		}

		if backupType == "One-time" {
			backupNextSnapshot = "N/A - One-time Backup"
		} else {
			backupNextSnapshot, _, err = unstructured.NestedString(backupJob.Object, "status", "nextSnapshotAt")
			if err != nil {
				return nil, err
			}
		}

		backupLastSnapshot, _, err = unstructured.NestedString(backupJob.Object, "status", "lastSnapshotTs")

		backupStorageLocation, _, err := unstructured.NestedString(backupJob.Object, "status", "storageLocation")
		if err != nil {
			return nil, err
		}

		backupFileName, _, err := unstructured.NestedString(backupJob.Object, "status", "filename")

		statusSlice, _, _ := unstructured.NestedSlice(backupJob.Object, "status", "conditions")
		for _, value := range statusSlice {
			for k, v := range value.(map[string]interface{}) {
				if k == "message" {
					backupMessage = v.(string)
				}
			}
		}

		backupInfo := backup{
			Name:            backupName,
			ResourceSetName: backupResourceSetName,
			RetentionCount:  backupRetentionCount,
			BackupType:      backupType,
			Message:         backupMessage,
			Filename:        backupFileName,
			NextSnapshot:    backupNextSnapshot,
			LastSnapshot:    backupLastSnapshot,
			StorageLocation: backupStorageLocation,
		}

		backups = append(backups, backupInfo)
	}
	return backups, nil
}

func (r Client) GetRestores() ([]restore, error) {
	var restores []restore
	var restoreMessage string

	res, err := r.Client.Resource(restoresGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, restoreJob := range res.Items {

		restoreName, _, err := unstructured.NestedString(restoreJob.Object, "metadata", "name")
		if err != nil {
			return nil, err
		}

		fileName, _, err := unstructured.NestedString(restoreJob.Object, "spec", "backupFilename")
		if err != nil {
			return nil, err
		}

		prune, _, err := unstructured.NestedBool(restoreJob.Object, "spec", "prune")
		if err != nil {
			return nil, err
		}

		restoreStorageLocation, _, err := unstructured.NestedString(restoreJob.Object, "status", "backupSource")
		if err != nil {
			return nil, err
		}

		statusSlice, _, err := unstructured.NestedSlice(restoreJob.Object, "status", "conditions")
		for _, value := range statusSlice {
			for k, v := range value.(map[string]interface{}) {
				if k == "message" {
					restoreMessage = v.(string)
				}
			}
		}
		if err != nil {
			return nil, err
		}

		restoreTime, _, err := unstructured.NestedString(restoreJob.Object, "status", "restoreCompletionTs")
		if err != nil {
			return nil, err
		}

		restoreInfo := restore{
			Name:                 restoreName,
			Filename:             fileName,
			Prune:                prune,
			StorageLocation:      restoreStorageLocation,
			Message:              restoreMessage,
			ResoreCompletionTime: restoreTime,
		}
		restores = append(restores, restoreInfo)
	}
	return restores, nil
}
