package handlers

import (
	"NetManager/env"
	"NetManager/logger"
	"NetManager/mqtt"
	"fmt"
	"net/http"
	"strings"
	"sync"
)

type ContainerDeployTask struct {
	Pid              int    `json:"pid"`
	NetworkNamespace string `json:"networkNamespace"`
	ServiceName      string `json:"serviceName"`
	Instancenumber   int    `json:"instanceNumber"`
	PortMappings     string `json:"portMappings"`
	Runtime          string
	PublicAddr       string
	PublicPort       string
	Env              *env.Environment
	Writer           *http.ResponseWriter
	Finish           chan TaskReady
}

type k8sDeployTask struct {
	Pid              int    `json:"pid"`
	NetworkNamespace string `json:"networkNamespace"`
	ServiceName      string `json:"serviceName"`
	Instancenumber   int    `json:"instanceNumber"`
	Podname          string `json:"podName"`
}

type TaskReady struct {
	deployment DeployResponse
	Err  error
}

type deployTaskQueue struct {
	newTask chan *ContainerDeployTask
}

type DeployTaskQueue interface {
	NewTask(request *ContainerDeployTask)
}

var (
	once      sync.Once
	taskQueue deployTaskQueue
)

func NewDeployTaskQueue() DeployTaskQueue {
	once.Do(func() {
		taskQueue = deployTaskQueue{
			newTask: make(chan *ContainerDeployTask, 50),
		}
		go taskQueue.taskExecutor()
	})
	return &taskQueue
}

func (t *deployTaskQueue) NewTask(request *ContainerDeployTask) {
	t.newTask <- request
}

func (t *deployTaskQueue) taskExecutor() {
	for {
		select {
		case task := <-t.newTask:
			// deploy the network stack in the container
			response, err := deploymentHandler(task)
			if err != nil {
				logger.ErrorLogger().Println("[ERROR]: ", err)
			}
			task.Finish <- TaskReady{
				deployment: *response,
				Err:  err,
			}
			// asynchronously update proxy tables
			updateInternalProxyDataStructures(task)
		}
	}
}

func deploymentHandler(requestStruct *ContainerDeployTask) (*DeployResponse, error) {
	// get app full name
	appCompleteName := strings.Split(requestStruct.ServiceName, ".")
	if len(appCompleteName) != 4 {
		return nil, fmt.Errorf("invalid app name: %s", appCompleteName)
	}

	// attach network to the container
	netHandler := env.GetNetDeployment(requestStruct.Runtime)
	logger.DebugLogger().Printf("Got netHandler: %v", netHandler)

	params, err := netHandler.DeployNetwork(requestStruct.Pid, requestStruct.NetworkNamespace, requestStruct.ServiceName, requestStruct.Instancenumber, requestStruct.PortMappings)
	logger.InfoLogger().Printf("Deployment: %+v", params)
	if err != nil {
		logger.ErrorLogger().Println("[ERROR]:", err)
		return nil, err
	}

	// notify to net-component
	err = mqtt.NotifyDeploymentStatus(
		requestStruct.ServiceName,
		"DEPLOYED",
		requestStruct.Instancenumber,
		params.ContainerIP.String(),
		params.ContainerIPv6.String(),
		requestStruct.PublicAddr,
		requestStruct.PublicPort,
	)

	if err != nil {
		logger.ErrorLogger().Println("[ERROR]:", err)
		return nil, err
	}

	response := &DeployResponse{
		ServiceName: params.ServiceName,
		HostVethName: params.HostVethName,
		HostBridgeName: params.HostBridgeName,
		HostBridgeIP: params.HostBridgeIP.String(),
		HostBridgeIPMask: params.HostBridgeIPMask,
		HostBridgeIPv6: params.HostBridgeIPv6.String(),
		HostBridgeIPv6Mask: params.HostBridgeIPv6Mask,
		ContainerVethName: params.ContainerVethName,
		ContainerNetNs: requestStruct.NetworkNamespace,
		ContainerIP: params.ContainerIP.String(),
		ContainerIPv6: params.ContainerIPv6.String(),
		Mtu: params.Mtu,
	}

	return response, nil
}

func updateInternalProxyDataStructures(requestStruct *ContainerDeployTask) {
	// Update internal table entry if an interest has not been set already.
	// Otherwise, do nothing, the net will autonomously update.
	if !mqtt.MqttIsInterestRegistered(requestStruct.ServiceName) {
		requestStruct.Env.RefreshServiceTable(requestStruct.ServiceName)
		mqtt.MqttRegisterInterest(requestStruct.ServiceName, requestStruct.Env, requestStruct.Instancenumber)
	}
}
