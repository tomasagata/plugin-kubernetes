package main

import (
	"log"
	"time"

	hardwareAgreggation "oakestra/plugin-kubernetes/oakestra-agent/agent/aggregation"
	config "oakestra/plugin-kubernetes/oakestra-agent/agent/config"
	kubenetesclient "oakestra/plugin-kubernetes/oakestra-agent/agent/kubernetesClient"
	kubenetesProxy "oakestra/plugin-kubernetes/oakestra-agent/agent/proxy"
)

func main() {
	config := config.GetConfig()

	k8sClient, err := kubenetesclient.NewKubernetesClient()
	if err != nil {
		log.Fatalln("Could not create Kubernetes Client:", err)
	}

	server := kubenetesProxy.NewServer(k8sClient, config)
	go server.StartHTTPServer()

	myClusterID, err := kubenetesProxy.Register(config)
	if err != nil {
		log.Fatalln("Could not register at root:", err)
	}

	config.SetClusterID(myClusterID)
	err = k8sClient.UpdateClusterInfoConfigMap(myClusterID, config)
	if err != nil {
		log.Fatalln("Could not update ConfigMap:", err)
	}

	timeInterval := 5 * time.Second
	err = hardwareAgreggation.SetUp()
	if err != nil {
		log.Fatalln("Could not create kubernetes Clients for Hardware Aggregation:", err)
	}
	go hardwareAgreggation.StartBackgroundServiceClusterInformation(timeInterval, config, k8sClient)
	select {}
}
