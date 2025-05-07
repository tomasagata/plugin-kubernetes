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

	// netnsPath := filepath.Join("/var/run/netns", netns)
	// fd, err := unix.Open(netnsPath, unix.O_RDONLY|unix.O_CLOEXEC, 0)
	// if err != nil {
	// 	log.Printf("File Descripter Error: %v", err)
	// }

	env := h.env
	// cleanup := func(veth *netlink.Veth) {
	// 	_ = netlink.LinkDel(veth)
	// }

	hostVethName, containerVethName := env.createVethsPairAndAttachToBridge(sname, env.mtusize)
	// if err != nil {
	// 	go cleanup(vethIfce)
	// 	return nil, nil, err
	// }

	// // Attach veth2 to the docker container
	// logger.DebugLogger().Println("Attaching peerveth to container ")
	// peerVeth, err := netlink.LinkByName(vethIfce.PeerName)
	// if err != nil {
	// 	cleanup(vethIfce)
	// 	return nil, nil, err
	// }

	// if pid == 0 {
	// 	if err := netlink.LinkSetNsFd(peerVeth, fd); err != nil {
	// 		cleanup(vethIfce)
	// 		return nil, nil, err
	// 	}
	// } else {
	// 	if err := netlink.LinkSetNsPid(peerVeth, pid); err != nil {
	// 		cleanup(vethIfce)
	// 		return nil, nil, err
	// 	}
	// }

	//generate a new ip for this container
	containerIP, err := env.generateAddress()
	if err != nil {
		// cleanup(vethIfce)
		return nil, err
	}

	// generate a new ipv6 for this container
	containerIPv6, err := env.generateIPv6Address()
	if err != nil {
		// cleanup(vethIfce)
		env.freeContainerAddress(containerIP)
		return nil, err
	}

	// set ip to the container veth
	// logger.DebugLogger().Println("Assigning ip ", ip.String()+env.config.HostBridgeMask, " to container ")
	// if err := env.addPeerLinkNetwork(pid, netnsPath, ip.String()+env.config.HostBridgeMask, vethIfce.PeerName); err != nil {
	// 	logger.ErrorLogger().Println("Error in addPeerLinkNetwork")
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return nil, nil, err
	// }

	// logger.DebugLogger().Println("Disabling DAD for IPv6")
	// if err := env.disableDAD(pid, netnsPath, vethIfce.PeerName); err != nil {
	// 	logger.ErrorLogger().Println("Error in Disabling DAD")
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return nil, nil, err
	// }

	// logger.DebugLogger().Println("Assigning ipv6 ", ipv6.String()+env.config.HostBridgeIPv6Prefix, " to container ")

	// if err := env.addPeerLinkNetwork(pid, netnsPath, ipv6.String()+env.config.HostBridgeIPv6Prefix, vethIfce.PeerName); err != nil {
	// 	logger.ErrorLogger().Println("Error in addPeerLinkNetworkv6")
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return nil, nil, err
	// }

	// Add traffic route to bridge
	// logger.DebugLogger().Println("Setting container routes ")
	// if err = env.setContainerRoutes(pid, netnsPath, vethIfce.PeerName); err != nil {
	// 	logger.ErrorLogger().Println("Error in setContainerRoutes")
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return nil, nil, err
	// }

	// if err = env.setIPv6ContainerRoutes(pid, netnsPath, vethIfce.PeerName); err != nil {
	// 	logger.ErrorLogger().Println("Error in setIPv6ContainerRoutes")
	// 	cleanup(vethIfce)
	// 	env.freeContainerAddress(ip)
	// 	env.freeContainerAddress(ipv6)
	// 	return nil, nil, err
	// }

	env.BookVethNumber()
	if err = env.setVethFirewallRules(hostVethName); err != nil {
		logger.ErrorLogger().Println("Error in setFirewallRules")
		// cleanup(vethIfce)
		env.freeContainerAddress(containerIP)
		env.freeContainerAddress(containerIPv6)
		return nil, err
	}

	if err = network.ManageContainerPorts(containerIP, portmapping, network.OpenPorts); err != nil {
		logger.ErrorLogger().Println("Error in ManageContainerPorts v4")
		debug.PrintStack()
		// cleanup(vethIfce)
		env.freeContainerAddress(containerIP)
		env.freeContainerAddress(containerIPv6)
		return nil, err
	}

	if err = network.ManageContainerPorts(containerIPv6, portmapping, network.OpenPorts); err != nil {
		logger.ErrorLogger().Println("Error in ManageContainerPorts v6")
		debug.PrintStack()
		// cleanup(vethIfce)
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
		ServiceName:     sname,
		HostVethName:    hostVethName,
		HostBridgeName:  env.config.HostBridgeName,
		HostBridgeIP:    net.ParseIP(env.config.HostBridgeIP),
		HostBridgeIPMask: env.config.HostBridgeMask,
		HostBridgeIPv6:  net.ParseIP(env.config.HostBridgeIPv6),
		HostBridgeIPv6Mask: env.config.HostBridgeIPv6Prefix,
		ContainerVethName: containerVethName,
		ContainerIP:     containerIP,
		ContainerIPv6:   containerIPv6,
		Mtu:            env.mtusize,
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
			HostVethName:    veths[0],
		}
		return params, nil
	}
	return nil, fmt.Errorf("service %s not found", snameAndInstance)
}
