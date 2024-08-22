package rancher

import (
	"context"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime/schema"
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

			go func(rancherCustomResource unstructured.Unstructured) {
				wg.Add(1)
				defer wg.Done()
				m.Lock()
				resource, group, _ := strings.Cut(rancherCustomResource.GetName(), ".")
				version, _, err := unstructured.NestedSlice(rancherCustomResource.Object, "status", "storedVersions")
				if err != nil {
					log.Errorf("error retrieving version of Rancher CRD: %v", err)
					m.Unlock() // Ensure the lock is released before returning
                    return      // Exit the goroutine early
				}

				result, err := r.Client.Resource(schema.GroupVersionResource{
					Group:    group,
					Version:  version[0].(string),
					Resource: resource,
				}).List(context.Background(), v1.ListOptions{})

				if err != nil {
					log.Errorf("error retrieving count of Rancher CRD: %v,%s,%s,%s\n", err, group, version[0].(string), resource)
					m.Unlock() // Ensure the lock is released before returning
                    return      // Exit the goroutine early
				}
				rancherCustomResources[rancherCustomResource.GetName()] = len(result.Items)
				m.Unlock()
			}(customResource)
		}
	}
	wg.Wait()
	return rancherCustomResources, nil
}
