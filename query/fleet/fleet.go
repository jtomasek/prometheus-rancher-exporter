package fleet

import (
	"context"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type Client struct {
	Client dynamic.Interface
	Config *rest.Config
}

var (
	settingGVRNumberOfNodes = schema.GroupVersionResource{Group: "fleet.cattle.io", Version: "v1alpha1", Resource: "clustergroups"}
)

func (r Client) GetNumberOfClusterGroups() (int, error) {

	res, err := r.Client.Resource(settingGVRNumberOfNodes).List(context.Background(), v1.ListOptions{})
	if err != nil {
		return 0, err
	}

	return len(res.Items), nil
}
