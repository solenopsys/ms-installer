package main

import (
	"git.klogsolenopsys.org/zmq_connector"
	"k8s.io/client-go/rest"
	"os"
)

var restClient *rest.Config

var devMode = os.Getenv("developerMode") == "true"

const ConfigmapName = "helm-repositories"
const NameSpace = "default"

func main() {
	var err error
	restClient, err = getCubeConfig(devMode)
	if err != nil {
		klog.Error("error getting Kubernetes config:", err)
		os.Exit(1)
	}
	template := zmq_connector.HsTemplate{Pf: processingFunction()}
	template.Init()
}
