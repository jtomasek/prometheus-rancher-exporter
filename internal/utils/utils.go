package utils

import (
	"context"
	"github.com/david-vtuk/prometheus-rancher-exporter/query/rancher"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
)

var (
	appGVR = schema.GroupVersionResource{Group: "catalog.cattle.io", Version: "v1", Resource: "apps"}
)

// CheckInstalledRancherApps /*
func CheckInstalledRancherApps(r rancher.Client) (bool, bool, error) {

	var rancherInstalled = false
	var rancherBackupInstalled = false

	res, err := r.Client.Resource(appGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return false, false, err
	}

	for _, app := range res.Items {

		appName, _, err := unstructured.NestedString(app.Object, "metadata", "name")
		if err != nil {
			return false, false, err
		}

		appStatus, _, err := unstructured.NestedString(app.Object, "spec", "info", "status")
		if err != nil {
			return false, false, err
		}

		if appName == "rancher" && appStatus == "deployed" {
			rancherInstalled = true
		}

		if appName == "rancher-backup" && appStatus == "deployed" {
			rancherBackupInstalled = true
		}

	}

	return rancherInstalled, rancherBackupInstalled, nil
}
