package main

import (
	"NetManager/env"
	"NetManager/handlers"
	kubenetesclient "NetManager/kuberneteslient.go"
	"NetManager/logger"
	"NetManager/mqtt"
	"NetManager/network"
	"NetManager/proxy"
	"crypto/sha256"
	"encoding/hex"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/pkg/errors"
)

type undeployRequest struct {
	Servicename    string `json:"serviceName"`
	Instancenumber int    `json:"instanceNumber"`
}

type registerRequest struct {
	ClientID string `json:"client_id"`
}

type DeployResponse struct {
	ServiceName string `json:"serviceName"`
	NsAddress   string `json:"nsAddress"`
}

type netConfiguration struct {
	NodePublicAddress string
	NodePublicPort    string
	ClusterUrl        string
	ClusterMqttPort   string
}

func handleRequests(port int) {
	netRouter := mux.NewRouter().StrictSlash(true)

	handlers.RegisterAllManagers(&Env, &WorkerID, Configuration.NodePublicAddress, Configuration.NodePublicPort, netRouter)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), netRouter))
}

func handleRequestsCNIPlugin(port int) {
	netRouter := mux.NewRouter().StrictSlash(true)
	handlers.RegisterAllManagers(&Env, &WorkerID, Configuration.NodePublicAddress, Configuration.NodePublicPort, netRouter)

	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), netRouter))
}

var Env env.Environment
var Proxy proxy.GoProxyTunnel
var WorkerID string
var Configuration netConfiguration

/*
Automatic register in k8s cluster
*/
func automaticRegister() error {
	log.Println("Start automatic register")
	hostname, err := os.Hostname()
	if err != nil {
		return fmt.Errorf("Could not retrieve hostname: %v", err)
	}
	clientID := hostname + "-" + Configuration.NodePublicAddress

	hash := sha256.Sum256([]byte(clientID))
	clientIDHashed := hex.EncodeToString(hash[:])[:24]

	WorkerID = clientIDHashed
	clientID = clientIDHashed

	//initialize mqtt connection to the broker
	mqtt.InitNetMqttClient(clientID, Configuration.ClusterUrl, Configuration.ClusterMqttPort)

	//initialize the proxy tunnel
	Proxy = proxy.New()
	Proxy.Listen()

	//initialize the Env Manager
	Env = *env.NewEnvironmentClusterConfigured(Proxy.HostTUNDeviceName)

	Proxy.SetEnvironment(&Env)
	return nil
}

func getHostIP() {

}

func main() {
	localPort := flag.Int("p", 6000, "Default local port of the NetManager")
	debugMode := flag.Bool("D", false, "Debug mode, it enables debug-level logs")
	flag.Parse()

	c, _ := kubenetesclient.NewKubernetesClient()
	hostIP, err := c.GetHostIP() // hostIP := "192.168.123.196"
	if err != nil {
		log.Fatal(errors.Wrap(err, "failed to load hostip"))
	}

	Configuration.NodePublicAddress = hostIP
	Configuration.NodePublicPort = os.Getenv("NODE_PORT")
	Configuration.ClusterMqttPort = os.Getenv("MOSQUITTO_SVC_SERVICE_PORT")
	Configuration.ClusterUrl = os.Getenv("MOSQUITTO_SVC_SERVICE_HOST")

	if *debugMode {
		logger.SetDebugMode()
	}
	network.IptableFlushAll()

	log.Println("NetManager started. Start Registration of Node.")

	automaticRegister()

	// Start manager for listenints
	handleRequestsCNIPlugin(*localPort)
}
