package env

import (
	"NetManager/logger"
	"NetManager/mqtt"
	"NetManager/network"
	"fmt"
	"net"
	"runtime/debug"
	"strings"
	// "github.com/vishvananda/netlink"
)

type ContainerDeyplomentHandler struct {
	env *Environment
}

var containerHandler *ContainerDeyplomentHandler = nil

func GetContainerNetDeployment() *ContainerDeyplomentHandler {
	if containerHandler == nil {
		logger.ErrorLogger().Fatal("Container Handler not initialized")
	}
	return containerHandler
}

func InitContainerDeployment(env *Environment) {
	containerHandler = &ContainerDeyplomentHandler{
		env: env,
	}
}

// AttachNetworkToContainer Attach a Docker container to the bridge and the current network environment
func (h *ContainerDeyplomentHandler) DeployNetwork(pid int, netns string, sname string, instancenumber int, portmapping string) (*DeploymentParameters, error) {
	env := h.env
	hostVethName, containerVethName := env.CreateVethsPair(sname)

	//generate a new ip for this container
	containerIP, err := env.generateAddress()
	if err != nil {
		return nil, err
	}

	// generate a new ipv6 for this container
	containerIPv6, err := env.generateIPv6Address()
	if err != nil {
		env.freeContainerAddress(containerIP)
		return nil, err
	}

	env.BookVethNumber()
	if err = env.setVethFirewallRules(hostVethName); err != nil {
		logger.ErrorLogger().Println("Error in setFirewallRules")
		env.freeContainerAddress(containerIP)
		env.freeContainerAddress(containerIPv6)
		return nil, err
	}

	if err = network.ManageContainerPorts(containerIP, portmapping, network.OpenPorts); err != nil {
		logger.ErrorLogger().Println("Error in ManageContainerPorts v4")
		debug.PrintStack()
		env.freeContainerAddress(containerIP)
		env.freeContainerAddress(containerIPv6)
		return nil, err
	}

	if err = network.ManageContainerPorts(containerIPv6, portmapping, network.OpenPorts); err != nil {
		logger.ErrorLogger().Println("Error in ManageContainerPorts v6")
		debug.PrintStack()
		env.freeContainerAddress(containerIP)
		env.freeContainerAddress(containerIPv6)
		return nil, err
	}

	env.deployedServicesLock.Lock()
	env.deployedServices[fmt.Sprintf("%s.%d", sname, instancenumber)] = service{
		ip:          containerIP,
		ipv6:        containerIPv6,
		sname:       sname,
		portmapping: portmapping,
		veth:        fmt.Sprintf("%s,%s", hostVethName, containerVethName),
	}
	env.deployedServicesLock.Unlock()
	logger.DebugLogger().Printf("New deployedServices table: %v", env.deployedServices)

	parameters := &DeploymentParameters{
		ServiceName:        sname,
		HostVethName:       hostVethName,
		HostBridgeName:     env.config.HostBridgeName,
		HostBridgeIP:       net.ParseIP(env.config.HostBridgeIP),
		HostBridgeIPMask:   env.config.HostBridgeMask,
		HostBridgeIPv6:     net.ParseIP(env.config.HostBridgeIPv6),
		HostBridgeIPv6Mask: env.config.HostBridgeIPv6Prefix,
		ContainerVethName:  containerVethName,
		ContainerIP:        containerIP,
		ContainerIPv6:      containerIPv6,
		Mtu:                env.mtusize,
	}
	return parameters, nil
}

func (env *Environment) DetachContainer(sname string, instance int) (*UndeploymentParameters, error) {
	snameAndInstance := fmt.Sprintf("%s.%d", sname, instance)
	env.deployedServicesLock.RLock()
	s, ok := env.deployedServices[snameAndInstance]
	env.deployedServicesLock.RUnlock()
	if ok {
		_ = env.translationTable.RemoveByNsip(s.ip)
		env.deployedServicesLock.Lock()
		delete(env.deployedServices, snameAndInstance)
		env.deployedServicesLock.Unlock()
		env.freeContainerAddress(s.ip)
		env.freeContainerAddress(s.ipv6)
		_ = network.ManageContainerPorts(s.ip, s.portmapping, network.ClosePorts)
		_ = network.ManageContainerPorts(s.ipv6, s.portmapping, network.ClosePorts)
		// _ = netlink.LinkDel(s.veth)
		// if no interest registered delete all remaining info about the service
		if !mqtt.MqttIsInterestRegistered(sname) {
			env.RemoveServiceEntries(sname)
		}

		veths := strings.Split(s.veth, ",")
		params := &UndeploymentParameters{
			HostVethName: veths[0],
		}
		return params, nil
	}
	return nil, fmt.Errorf("service %s not found", snameAndInstance)
}
