package main

import (
	"flag"
	"fmt"
	"github.com/david-vtuk/prometheus-rancher-exporter/internal/utils"
	"github.com/prometheus/client_golang/prometheus"
	"net/http"
	"os"
	"os/user"
	"strconv"
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

func getEnv(key string, defaultValue string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}

	return defaultValue
}

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

	Timer_GetLatestRancherVersion, err := strconv.Atoi(getEnv("TIMER_GET_LATEST_RANCHER_VERSION", "1"))
	Timer_ticker, err := strconv.Atoi(getEnv("TIMER_TICKER", "10"))

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

	rancherInstalled, rancherBackupsInstalled, err := utils.CheckInstalledRancherApps(RancherClient)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	log.Printf("Detected Rancher: %s", strconv.FormatBool(rancherInstalled))
	log.Printf("Detected Rancher Backup Operator: %s", strconv.FormatBool(rancherBackupsInstalled))

	//Kick off collector in background
	if rancherInstalled {
		log.Printf("Collecting Rancher Metrics")
		http.Handle("/metrics", promhttp.Handler())
		go collector.Collect(RancherClient, Timer_GetLatestRancherVersion, Timer_ticker, rancherBackupsInstalled)
	}

	if rancherBackupsInstalled {
		log.Printf("Collecting Rancher Backup Operator Metrics")
		reg := prometheus.NewRegistry()
		backupsHandler := promhttp.HandlerFor(reg, promhttp.HandlerOpts{})
		http.Handle("/backup-metrics", backupsHandler)
		go collector.CollectBackupMetrics(RancherClient, Timer_ticker, reg)
	}

	log.Info("Beginning to serve on port :8080")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
