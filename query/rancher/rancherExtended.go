package rancher

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"time"
)

var (
	crdGVR               = schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	customResourceDomain = "cattle.io"
)

func (r Client) GetRancherCustomResourceCount() (map[string]int, error) {

	start := time.Now()
	rancherCustomResources := make(map[string]int)

	res, err := r.Client.Resource(crdGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, customResource := range res.Items {

		if strings.Contains(customResource.GetName(), customResourceDomain) {

			// put this in a goroutine, probably need a mutex.

			resource, group, _ := strings.Cut(customResource.GetName(), ".")
			version, _, err := unstructured.NestedSlice(customResource.Object, "status", "storedVersions")
			if err != nil {
				return nil, err
			}

			result, err := r.Client.Resource(schema.GroupVersionResource{
				Group:    group,
				Version:  version[0].(string),
				Resource: resource,
			}).List(context.Background(), v1.ListOptions{})

			if err != nil {
				return nil, err
			}

			rancherCustomResources[customResource.GetName()] = len(result.Items)
		}
	}

	elapsed := time.Since(start)
	fmt.Printf("CRD took %s", elapsed)

	return rancherCustomResources, nil
}
