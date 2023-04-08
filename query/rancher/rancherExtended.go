package rancher

import (
	"context"
	"fmt"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"strings"
)

var (
	crdGVR               = schema.GroupVersionResource{Group: "apiextensions", Version: "v1", Resource: "CustomResourceDefinition"}
	customResourceDomain = "cattle.io"
)

func (r Client) GetRancherCustomResourceCount() (map[string]int, error) {

	res, err := r.Client.Resource(crdGVR).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return nil, err
	}

	for _, customResource := range res.Items {

		if strings.Contains(customResource.GetName(), customResourceDomain) {

			fmt.Println(customResource.GetName())
		}
	}
	return nil, nil
}
