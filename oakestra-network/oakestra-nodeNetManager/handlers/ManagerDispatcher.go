package handlers

import (
	"NetManager/env"
	"log"

	"github.com/gorilla/mux"
)

type netConfiguration struct {
	NodePublicAddress string
	NodePublicPort    string
}

type undeployRequest struct {
	Servicename    string `json:"serviceName"`
	Instancenumber int    `json:"instanceNumber"`
}

type undeployResponse struct {
	VethPeer1Name string `json:"vethPeer1Name"`
}

type DeployResponse struct {
	ServiceName string `json:"serviceName"`
	HostVethName string `json:"hostVethName"`
	HostBridgeName string `json:"hostBridgeName"`
	HostBridgeIP string `json:"hostBridgeIP"`
	HostBridgeIPMask string `json:"hostBridgeIPMask"`
	HostBridgeIPv6 string `json:"hostBridgeIPv6"`
	HostBridgeIPv6Mask string `json:"hostBridgeIPv6Mask"`
	ContainerVethName string `json:"containerVethName"`
	ContainerNetNs string `json:"containerNetNs"`
	ContainerIP string `json:"containerIP"`
	ContainerIPv6 string `json:"containerIPv6"`
	Mtu       int    `json:"mtu"`
}

var AvailableRuntimes = make(map[string]func() ManagerInterface)

type ManagerInterface interface {
	Register(Env *env.Environment, WorkerID *string, NodePublicAddress string, NodePublicPort string, Router *mux.Router)
}

func GetNetManager(handler string) ManagerInterface {
	if getfunc, ok := AvailableRuntimes[handler]; ok {
		return getfunc()
	}
	return nil
}

func RegisterAllManagers(Env *env.Environment, WorkerID *string, NodePublicAddress string, NodePublicPort string, Router *mux.Router) {

	for s, getfunc := range AvailableRuntimes {

		log.Println("Available Runtime")
		log.Println(s)
		getfunc().Register(Env, WorkerID, NodePublicAddress, NodePublicPort, Router)
	}
}
