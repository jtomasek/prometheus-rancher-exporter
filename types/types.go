package types

import (
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
)

type Client struct {
	Client dynamic.Interface
	Config *rest.Config
}
