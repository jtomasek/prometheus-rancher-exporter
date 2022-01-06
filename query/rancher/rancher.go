package rancher

import (
	"context"
	"github.com/ebauman/prometheus-rancher-exporter/semver"
	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
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

	log.Info("result: ", res)

	log.Info("unstructured content ", res.UnstructuredContent())

	version, _, _ := unstructured.NestedString(res.UnstructuredContent(), "value")

	log.Info("version", version)

	// major, minor, patch, prerelease, buildmetadata
	result, err := semver.Parse("v2.6.3")
	if err != nil {
		return nil, err
	}

	return result, nil
}
