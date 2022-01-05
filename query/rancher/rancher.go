package rancher

import (
	"context"
	"github.com/ebauman/prometheus-rancher-exporter/semver"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

var (
	settingGVR = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "settings"}
)

type Client struct {
	client dynamic.Interface
}

/*
func (r Client) NewClient() (dynamic.Interface, error) {
	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}
	r.client, err = dynamic.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	return r.client, err
}
*/

func (r Client) GetRancherVersion() (map[string]int64, error) {

	config, err := rest.InClusterConfig()

	if err != nil {
		return nil, err
	}
	r.client, err = dynamic.NewForConfig(config)

	if err != nil {
		return nil, err
	}

	res, err := r.client.Resource(settingGVR).Get(context.Background(), "server-version", v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	// major, minor, patch, prerelease, buildmetadata
	result, err := semver.Parse(res.GetName())
	if err != nil {
		return nil, err
	}

	return result, nil
}
