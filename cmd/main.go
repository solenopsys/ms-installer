package main

import (
	zmq_connector "github.com/solenopsys/sc-bl-zmq-connector"
	"hs-installer/internal"
	"k8s.io/klog/v2"
	"os"
)

var devMode = os.Getenv("developerMode") == "true"

const ConfigmapName = "helm-repositories"
const NameSpace = "default"

func main() {
	var err error

	if err != nil {
		klog.Error("error getting Kubernetes config:", err)
		os.Exit(1)
	}
	template := zmq_connector.HsTemplate{Pf: internal.ProcessingFunction()}
	template.Init()
}
