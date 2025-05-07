package env

import "net"

const (
	CONTAINER_RUNTIME = "container"
)

type NetDeploymentInterface interface {
	DeployNetwork(pid int, netns string, sname string, instancenumber int, portmapping string) (*DeploymentParameters, error)
}

func GetNetDeployment(handler string) NetDeploymentInterface {
	switch handler {
	case CONTAINER_RUNTIME:
		return GetContainerNetDeployment()
	}
	return nil
}

type DeploymentParameters struct {
	ServiceName string
	HostVethName string
	HostBridgeName string
	HostBridgeIP net.IP
	HostBridgeIPMask string
	HostBridgeIPv6 net.IP
	HostBridgeIPv6Mask string
	ContainerVethName string
	ContainerIP net.IP
	ContainerIPv6 net.IP
	Mtu int
}

type UndeploymentParameters struct {
	HostVethName string
}