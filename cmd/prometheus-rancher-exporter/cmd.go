package main

import (
	"flag"
	"github.com/ebauman/prometheus-rancher-exporter/collector"
)

func main() {


	go collector.Collect()
}