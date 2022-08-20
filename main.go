package main

import (
	"k8s.io/client-go/kubernetes"
	"os"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"solenopsys.org/alexstorm/zmq_connector"
)

var clientSet *kubernetes.Clientset
var c client.Client
var devMode = os.Getenv("developerMode") == "true"

const ConfigmapName = "helm-repositories"
const NameSpace = "default"

func main() {
	clientSet, c = createKubeConfig()
	template := zmq_connector.HsTemplate{Pf: processingFunction()}
	template.Init()
}
