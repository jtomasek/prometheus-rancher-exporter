package main

import (
	"github.com/ebauman/prometheus-rancher-exporter/collector"
	"github.com/ebauman/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"net/http"
)

func main() {

	// Build Rancher Client
	log.Info("Building Rancher Client")
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatal("Unable to construct Rancher client Config")
	}

	client, err := dynamic.NewForConfig(config)
	if err != nil {
		log.Fatal("Unable to construct Rancher client")
	}

	RancherClient := rancher.Client{
		Config: config,
		Client: client,
	}

	//Kick off collector in background
	go collector.Collect(RancherClient)

	//This section will start the HTTP server and expose
	//any metrics on the /metrics endpoint.
	http.Handle("/metrics", promhttp.Handler())
	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
