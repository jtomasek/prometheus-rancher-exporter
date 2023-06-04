package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/user"
	"strings"

	"github.com/david-vtuk/prometheus-rancher-exporter/collector"
	"github.com/david-vtuk/prometheus-rancher-exporter/query/rancher"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

const (
	k8sClientBurst = 100
	k8sClientQPS   = 100
)

func main() {

	// Build Rancher Client
	log.Info("Building Rancher Client")

	// Default to in-cluster config
	InClusterConfig := true

	var config *rest.Config
	var err error

	if strings.ToUpper(os.Getenv("RANCHER_EXPORTER_EXTERNAL_AUTH")) == "TRUE" {
		log.Info("RANCHER_EXPORTER_EXTERNAL_AUTH env variable set to true, using out of cluster config")
		InClusterConfig = false
	}

	if InClusterConfig {
		config, err = rest.InClusterConfig()
                if err != nil {
                        log.Fatal("Unable to construct REST client")
                }

		config.Burst = k8sClientBurst
		config.QPS = k8sClientQPS
	} else {
		currentUser, err := user.Current()
		if err != nil {
			log.Fatal(err.Error())
		}

		kubeconfig := flag.String("kubeconfig", fmt.Sprintf("/home/%s/.kube/config", currentUser.Username), "absolute path to the kubeconfig file")
		flag.Parse()
		config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
                if err != nil {
                        log.Fatal("Unable to construct Rancher client Config")
                }

		config.Burst = k8sClientBurst
		config.QPS = k8sClientQPS

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
