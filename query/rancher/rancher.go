package rancher

import (
	"context"
	"github.com/ebauman/prometheus-rancher-exporter/semver"
	"github.com/tidwall/gjson"
	"io"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"net/http"
)

var (
	settingGVR                 = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "settings"}
	settingGVRNumberOfClusters = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "clusters"}
	settingGVRNumberOfNodes    = schema.GroupVersionResource{Group: "management.cattle.io", Version: "v3", Resource: "nodes"}
)

type Client struct {
	Client dynamic.Interface
	Config *rest.Config
}

func (r Client) GetRancherVersion() (map[string]int64, error) {

	res, err := r.Client.Resource(settingGVR).Get(context.Background(), "server-version", v1.GetOptions{})
	if err != nil {
		return nil, err
	}

	version, _, _ := unstructured.NestedString(res.UnstructuredContent(), "value")

	result, err := semver.Parse(TrimVersionChar(version))

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Client) GetNumberOfManagedClusters() (int, error) {

	res, err := r.Client.Resource(settingGVRNumberOfClusters).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

func (r Client) GetK8sDistributions() (map[string]int, error) {

	distributions := make(map[string]int)

	res, err := r.Client.Resource(settingGVRNumberOfClusters).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, v := range res.Items {
		labels := v.GetLabels()
		distribution := labels["provider.cattle.io"]
		distributions[distribution] += 1
	}

	return distributions, nil
}

func (r Client) GetLatestRancherVersion() (map[string]int64, error) {
	resp, err := http.Get("https://api.github.com/repos/rancher/rancher/releases/latest")
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}

	bodyBytes, err := io.ReadAll(resp.Body)

	if err != nil {
		return nil, err
	}

	val := gjson.Get(string(bodyBytes), "tag_name")

	result, err := semver.Parse(TrimVersionChar(val.String()))

	if err != nil {
		return nil, err
	}

	return result, nil
}

func (r Client) GetNumberOfManagedNodes() (int, error) {

	res, err := r.Client.Resource(settingGVRNumberOfNodes).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}

// Version returned from CRD is in the format of "vN.N.N", trim the leading "v"
func TrimVersionChar(version string) string {
	for i := range version {
		if i > 0 {
			return version[i:]
		}
	}
	return version[:0]
}
