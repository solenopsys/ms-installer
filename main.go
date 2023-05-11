package main

import (
	bl_kubernetes_tools "github.com/solenopsys/bl-kubernetes-tools"
	zmq_connector "github.com/solenopsys/sc-bl-zmq-connector"
	"k8s.io/client-go/rest"
	"k8s.io/klog/v2"
	"os"
)

var restClient *rest.Config

var devMode = os.Getenv("developerMode") == "true"

const ConfigmapName = "helm-repositories"
const NameSpace = "default"

func main() {
	var err error
	restClient, err = bl_kubernetes_tools.GetCubeConfig(devMode)
	if err != nil {
		klog.Error("error getting Kubernetes config:", err)
		os.Exit(1)
	}
	template := zmq_connector.HsTemplate{Pf: processingFunction()}
	template.Init()
}
