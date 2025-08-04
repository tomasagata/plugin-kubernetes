package env

import (
	"oakestra/cni-plugin/logger"
	"oakestra/cni-plugin/models"
	"net"
	"os/exec"
	"runtime"

	"github.com/vishvananda/netlink"
	"github.com/vishvananda/netns"
)

// create veth pair and connect one to the host bridge
// returns: bridgeVeth name, free Veth name, Vether interface to the veth pair and eventually an error
func createVethsPairAndAttachToBridge(details *models.ConnectNetworkResponse) (*netlink.Veth, error) {
	// Retrieve current bridge
	logger.DebugLogger().Println("Retrieving current bridge ")
	bridge, err := netlink.LinkByName(details.HostBridgeName)
	if err != nil {
		logger.ErrorLogger().Println("Error retrieving current bridge: ", err)
		return nil, err
	}
	logger.DebugLogger().Println("Retrieved current bridge")
	logger.InfoLogger().Println("Bridge: " + bridge.Attrs().Name)
	logger.DebugLogger().Println("creating veth pair: " + details.HostVethName + "@" + details.ContainerVethName)

	veth := &netlink.Veth{
		LinkAttrs: netlink.LinkAttrs{
			Name: details.HostVethName,
			MTU:  details.Mtu,
		},
		PeerName: details.ContainerVethName,
	}
	err = netlink.LinkAdd(veth)
	if err != nil {
		return nil, err
	}

	// add veth1 to the bridge
	err = netlink.LinkSetMaster(veth, bridge)
	if err != nil {
		return nil, err
	}

	// set veth status up
	if err = netlink.LinkSetUp(veth); err != nil {
		return nil, err
	}

	return veth, nil
}

// add routes inside the container namespace to forward the traffic using the bridge
func setContainerRoutes(netnsPath string, peerVeth string, hostBridgeIP string) error {
	//Add route to bridge
	//sudo nsenter -n -t 5565 ip route add 0.0.0.0/0 via 127.19.x.y dev veth013
	err := execInsideNs(netnsPath, func() error {
		link, err := netlink.LinkByName(peerVeth)
		if err != nil {
			return err
		}
		dst, err := netlink.ParseIPNet("10.30.0.0/16")
		if err != nil {
			return err
		}
		gw := net.ParseIP(hostBridgeIP)
		return netlink.RouteAdd(&netlink.Route{
			LinkIndex: link.Attrs().Index,
			Dst:       dst,
			Gw:        gw,
		})
	})
	if err != nil {
		logger.ErrorLogger().Printf("Impossible to setup route inside the netns: %v\n", err)
		return err
	}
	return nil
}

func setIPv6ContainerRoutes(netnsPath string, peerVeth string, hostBridgeIPv6 string) error {
	err := execInsideNs(netnsPath, func() error {
		link, err := netlink.LinkByName(peerVeth)
		if err != nil {
			return err
		}
		dstv6, err := netlink.ParseIPNet("fdff:2000::/21")
		if err != nil {
			return err
		}
		gwv6 := net.ParseIP(hostBridgeIPv6)
		return netlink.RouteAdd(&netlink.Route{
			LinkIndex: link.Attrs().Index,
			Dst:       dstv6,
			Gw:        gwv6,
		})
	})
	if err != nil {
		logger.ErrorLogger().Printf("Impossible to setup IPv6 route inside the netns: %v\n", err)
		return err
	}
	return nil
}

// setup the address of the network namespace veth
func addPeerLinkNetwork(netnsPath string, addr string, vethname string) error {
	netlinkAddr, err := netlink.ParseAddr(addr)
	if err != nil {
		return err
	}
	err = execInsideNs(netnsPath, func() error {
		link, err := netlink.LinkByName(vethname)
		if err != nil {
			return err
		}
		err = netlink.AddrAdd(link, netlinkAddr)
		if err == nil {
			err = netlink.LinkSetUp(link)
		}
		return err
	})
	if err != nil {
		return err
	}
	return err
}

// setup the address of the network namespace veth based on Ns name
func addPeerLinkNetworkByNsName(NsName string, addr string, vethname string) error {
	netlinkAddr, err := netlink.ParseAddr(addr)
	if err != nil {
		return err
	}
	err = execInsideNsByName(NsName, func() error {
		link, err := netlink.LinkByName(vethname)
		if err != nil {
			return err
		}
		err = netlink.AddrAdd(link, netlinkAddr)
		if err == nil {
			err = netlink.LinkSetUp(link)
		}
		return err
	})
	return err
}

// disable Duplicate Address Detection (DAD) for IPv6 interfaces in namespace
// to prevent interface startup delay
func disableDAD(netnsPath string, vethname string) error {
	err := execInsideNs(netnsPath, func() error {
		cmd := exec.Command("sysctl", "-w", "net.ipv6.conf.default.accept_dad=0")
		err := cmd.Run()
		if err != nil {
			return err
		}
		cmd = exec.Command("sysctl", "-w", "net.ipv6.conf."+vethname+".accept_dad=0")
		err = cmd.Run()
		return err
	})
	return err
}

// Execute function inside a namespace
func execInsideNs(netnsPath string, function func() error) error {
	var containerNs netns.NsHandle

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	stdNetns, err := netns.Get()
	if err == nil {
		defer stdNetns.Close()
		containerNs, err = netns.GetFromPath(netnsPath)
		if err == nil {
			defer netns.Set(stdNetns)
			err = netns.Set(containerNs)
			if err == nil {
				err = function()
			}
		}
	}
	return err
}

// Execute function inside a namespace based on Ns name
func execInsideNsByName(Nsname string, function func() error) error {
	var containerNs netns.NsHandle

	runtime.LockOSThread()
	defer runtime.UnlockOSThread()

	stdNetns, err := netns.Get()
	if err == nil {
		defer stdNetns.Close()
		containerNs, err = netns.GetFromName(Nsname)
		if err == nil {
			defer netns.Set(stdNetns)
			err = netns.Set(containerNs)
			if err == nil {
				err = function()
			}
		}
	}
	return err
}
