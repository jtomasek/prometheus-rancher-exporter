package rancher

import (
	"context"
	"github.com/prometheus/common/log"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
	"sync"
)

var (
	crdGVR               = schema.GroupVersionResource{Group: "apiextensions.k8s.io", Version: "v1", Resource: "customresourcedefinitions"}
	customResourceDomain = "cattle.io"
)

func (r Client) GetRancherCustomResourceCount() (map[string]int, error) {

	rancherCustomResources := make(map[string]int)
	var m sync.Mutex
	var wg sync.WaitGroup

	res, err := r.Client.Resource(crdGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, customResource := range res.Items {

		if strings.Contains(customResource.GetName(), customResourceDomain) {
			wg.Add(1)
			go func(rancherCustomResource unstructured.Unstructured) {
				defer wg.Done()
				m.Lock()
				resource, group, _ := strings.Cut(rancherCustomResource.GetName(), ".")
				version, _, err := unstructured.NestedSlice(rancherCustomResource.Object, "status", "storedVersions")
				if err != nil {
					log.Errorf("error retrieving version of Rancher CRD: %v", err)
				}

				result, err := r.Client.Resource(schema.GroupVersionResource{
					Group:    group,
					Version:  version[0].(string),
					Resource: resource,
				}).List(context.Background(), v1.ListOptions{})

				if err != nil {
					log.Errorf("error retrieving count of Rancher CRD: %v,%s,%s,%s\n", err, group, version[0].(string), resource)
				}
				rancherCustomResources[rancherCustomResource.GetName()] = len(result.Items)
				m.Unlock()
			}(customResource)
		}
	}
	wg.Wait()
	return rancherCustomResources, nil
}
