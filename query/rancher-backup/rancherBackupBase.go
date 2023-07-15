package rancherbackup

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	backupsGVR     = schema.GroupVersionResource{Group: "resources.cattle.io", Version: "v1", Resource: "backups"}
	restoresGVR    = schema.GroupVersionResource{Group: "resources.cattle.io", Version: "v1", Resource: "restores"}
	restoresetsGVR = schema.GroupVersionResource{Group: "resources.cattle.io", Version: "v1", Resource: "restoresets"}
)

type Client struct {
	Client dynamic.Interface
	Config *rest.Config
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

func (r Client) GetNumberOfRestoreSets() (int, error) {

	res, err := r.Client.Resource(restoresetsGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}
