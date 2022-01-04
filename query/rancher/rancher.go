package rancher

import (
	"context"
	"github.com/ebauman/prometheus-rancher-exporter/semver"
	managementClient "github.com/rancher/rancher/pkg/client/generated/management/v3"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
)

var (
	settingGVR = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "settings"}
)

type Client struct {
	client dynamic.Interface
}

func (r Client) GetRancherVersion() (map[string]int64, error) {
	res, err := r.client.Resource(settingGVR).Get(context.Background(), "server-version", v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	

	// major, minor, patch, prerelease, buildmetadata
	result, err := semver.Parse(setting.Value)
	if err != nil {
		return nil, err
	}

	return result, nil
}

