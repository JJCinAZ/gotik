package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"regexp"
	"strings"
	"time"

	"github.com/davecgh/go-spew/spew"

	"github.com/jjcinaz/gotik"
)

func main() {
	var (
		routerConn    *routeros.Client
		err           error
		ipsRemoved    []routeros.IPv4Address
		routesRemoves []routeros.IPv4Route
	)
	routerConn, err = routeros.DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"), time.Second*5)
	if err != nil {
		log.Fatal(err)
	} else {
		defer routerConn.Close()
		ipsRemoved, routesRemoves, err = practice(routerConn, 105, os.Getenv("RTR"), "sfp2-Backhaul-Sw1")
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println("success")
			spew.Dump(ipsRemoved)
			spew.Dump(routesRemoves)
		}
	}
}

/*
Cancel Service Algorithm:
	- DeprovisionRouterInterface()
	- Set SM.VAN = 4094
	- Send Pick-up Schedule Task email
*/

/*
DeprovisionRouterInterface clean up common provisioning for customer VLANs on routers
The following operations are executed:
	- Name and VLAN ID sanity check. Name must start with 'xxx-v###-' and VLAN ID must be 100...999
	- Check to see that we have a protected IP on the VLAN (172.16.0.0/12)
	- Removes any static routes (those removed are returned in routesRemoved slice)
	- Removes all IPv4 addresses which are not in 172.16.0.0/12 (those removed are returned in ipsRemoved slice)
	- Enable protected address if it was disabled
	- Enable DHCPv4 server for VLAN if it was disabled
	- Removes all queue tree entries for VLAN
	- Removes all simple queues for VLAN
	- Removes all comments from VLAN
	- Renames VLAN to xxx-v###-avail
	- Go through list of all removed IP addresses and remove any assignment in IPAM
	- Go through list of removed routes and remove those from any assignment in IPAM
*/
func DeprovisionRouterInterface(conn *routeros.Client, vlanid int, routerIP string, ringIntfName string) (
	ipsRemoved []routeros.IPv4Address, routesRemoved []routeros.IPv4Route, err error,
) {
	var (
		nameRegex1    = regexp.MustCompile(`(?i)(^\w+-v\d{1,4}-)`)
		nameParts     []string
		targIntf      routeros.Interface
		ipv4List      []routeros.IPv4Address
		addressToKeep *routeros.IPv4Address
		routeList     []routeros.IPv4Route
		pppSecrets    []routeros.PPPSecret
	)
	ipsRemoved = make([]routeros.IPv4Address, 0)
	routesRemoved = make([]routeros.IPv4Route, 0)
	_, protectedNet, _ := net.ParseCIDR("172.16.0.0/12")

	// Get Interface name for VLAN interface on ringIntfName
	targIntf, err = conn.GetVLANInterfaceOnBase(ringIntfName, vlanid)
	if err != nil {
		return
	}
	if targIntf.VLAN < 100 || targIntf.VLAN > 999 {
		err = fmt.Errorf("invalid VLAN ID range (should be 100...999) for vlan %d (%s)", vlanid, targIntf.Name)
		return
	}
	if nameParts = nameRegex1.FindStringSubmatch(targIntf.Name); nameParts == nil || len(nameParts) < 2 {
		err = fmt.Errorf("invalid name format for vlan %d (%s)", vlanid, targIntf.Name)
		return
	}

	// Find all IPv4 addresses on that VLAN
	// If the number of addresses found in 172.16.0.0/12 is not one then error
	ipv4List, err = conn.GetInterfaceIPv4Table(targIntf.Name)
	if err != nil {
		return
	}
	if ipv4List == nil || len(ipv4List) == 0 {
		err = fmt.Errorf("no IPv4 addresses found on vlan %d (%s)", vlanid, targIntf.Name)
		return
	}
	for i, addr := range ipv4List {
		if protectedNet.Contains(net.ParseIP(addr.Address)) {
			if addressToKeep != nil {
				err = fmt.Errorf("duplicate protected addresses on vlan %d (%s)", vlanid, targIntf.Name)
				return
			}
			addressToKeep = &ipv4List[i]
		}
	}
	if addressToKeep == nil {
		err = fmt.Errorf("missing address in the 172.16.0.0/12 subnet on vlan %d (%s)", vlanid, targIntf.Name)
		return
	}

	// 	Find and remove all IPv4 routes where static=true and where gateway-status ends with <TargIntf>
	if routeList, err = conn.GetIPv4Routes([]string{"static"}); err != nil {
		return
	}
	for _, r := range routeList {
		if strings.HasSuffix(r.GatewayStatus, targIntf.Name) {
			if err = conn.ModifyIPv4Route(r.ID, "remove"); err != nil {
				return
			}
			routesRemoved = append(routesRemoved, r)
		}
	}

	// If any removed routes were a /32 route, then see if that entry is in the PPP Secrets table
	if pppSecrets, err = conn.GetPPPSecrets(); err != nil {
		return
	}
	for _, s := range pppSecrets {
		if find32Route(routesRemoved, s.RemoteAddress) {
			s.Name = fmt.Sprintf("Static-%s", s.RemoteAddress)
			s.Comment = ""
			if _, err = conn.UpdatePPPSecret(s); err != nil {
				return
			}
			ipsRemoved = append(ipsRemoved, routeros.IPv4Address{
				ID:      "",
				Address: s.RemoteAddress + "/32",
				Network: s.RemoteAddress,
			})
		}
	}

	// Remove any PPPoE servers for TargIntf
	if pppoeServers, e := conn.GetPPPoEServers(targIntf.Name); e != nil {
		err = e
		return
	} else {
		for i := range pppoeServers {
			if err = conn.RemovePPPoEServer(pppoeServers[i].ID); err != nil {
				return
			}
		}
	}

	// Remove all IPv4 addresses except those in 172.16.0.0/12
	// Enable all remaining addresses (should just be the 172.16.0.0/12 one) if they were disabled
	for _, addr := range ipv4List {
		if addr.ID != addressToKeep.ID {
			if err = conn.ModifyIPv4Address(addr.ID, "remove"); err != nil {
				return
			}
			ipsRemoved = append(ipsRemoved, addr)
		} else if addr.Disabled {
			if err = conn.ModifyIPv4Address(addr.ID, "enable"); err != nil {
				return
			}
		}
	}

	// Find all queue tree entries where parent=<TargIntf>
	// Note that we do this here to allow some time delay after possibly enabling the IPv4 address and
	// enabling the DHCP server below.
	if queues, e := conn.GetQueueTree(targIntf.Name); e != nil {
		err = e
		return
	} else {
		for i := range queues {
			// We pass true here to recursively remove all child queues
			if err = conn.RemoveQueueTree(queues[i], true); err != nil {
				return
			}
		}
	}

	// Find all simply queue entries where target=<TargIntf>
	if simpQueues, e := conn.GetSimpleQueues(targIntf.Name); e != nil {
		err = e
		return
	} else {
		for _, q := range simpQueues {
			if err = conn.RemoveSimpleQueue(q.ID); err != nil {
				return
			}
		}
	}

	// Find DHCP server where interface=<TargIntf> and enable server if it was disabled
	if dhcpServers, e := conn.GetDhcp4ServerByIntf(targIntf.Name); e != nil {
		err = e
		return
	} else {
		if len(dhcpServers) > 1 {
			err = fmt.Errorf("more than one DHCPv4 server on vlan %d (%s)", vlanid, targIntf.Name)
			return
		}
		for _, s := range dhcpServers {
			if s.Disabled {
				if err = conn.SetDhcpv4ServerDisable(s.ID, false); err != nil {
					return
				}
			}
		}
	}

	// Remove comment from TargIntf
	if len(targIntf.Comment) > 0 {
		if err = conn.SetInterfaceComment(targIntf.ID, ""); err != nil {
			return
		}
	}

	// Rename TargIntf to `(^\w+-v\d{1,4}-)` + "avail"
	if err = conn.SetInterfaceName(targIntf.ID, nameParts[1]+"avail"); err != nil {
		return
	}
	return
}

func find32Route(routes []routeros.IPv4Route, address string) bool {
	if !strings.HasSuffix(address, "/32)") {
		address += "/32"
	}
	for _, r := range routes {
		if address == r.DstAddress {
			return true
		}
	}
	return false
}

func practice(conn *routeros.Client, vlanid int, routerIP string, ringIntfName string) (
	ipsRemoved []routeros.IPv4Address, routesRemoved []routeros.IPv4Route, err error,
) {
	var (
		nameRegex1    = regexp.MustCompile(`(?i)(^\w+-v\d{1,4}-)`)
		nameParts     []string
		targIntf      routeros.Interface
		ipv4List      []routeros.IPv4Address
		protectedNet  *net.IPNet
		addressToKeep *routeros.IPv4Address
		routeList     []routeros.IPv4Route
		pppSecrets    []routeros.PPPSecret
	)
	ipsRemoved = make([]routeros.IPv4Address, 0)
	routesRemoved = make([]routeros.IPv4Route, 0)

	// Get Interface name for VLAN interface on ringIntfName
	targIntf, err = conn.GetVLANInterfaceOnBase(ringIntfName, vlanid)
	if err != nil {
		return
	}
	if targIntf.VLAN < 100 || targIntf.VLAN > 999 {
		err = fmt.Errorf("invalid VLAN ID range (should be 100...999) for vlan %d (%s)", vlanid, targIntf.Name)
		return
	}
	if nameParts = nameRegex1.FindStringSubmatch(targIntf.Name); nameParts == nil || len(nameParts) < 2 {
		err = fmt.Errorf("invalid name format for vlan %d (%s)", vlanid, targIntf.Name)
		return
	}

	// Find all IPv4 addresses on that VLAN
	// If the number of addresses found in 172.16.0.0/12 is not one then error
	ipv4List, err = conn.GetInterfaceIPv4Table(targIntf.Name)
	if err != nil {
		return
	}
	if ipv4List == nil || len(ipv4List) == 0 {
		err = fmt.Errorf("no IPv4 addresses found on vlan %d (%s)", vlanid, targIntf.Name)
		return
	}
	if _, protectedNet, err = net.ParseCIDR("172.16.0.0/12"); err != nil {
		return
	}
	for i, addr := range ipv4List {
		// Note that we just check the Network here not the CIDR address
		// While we could parse the CIDR and then get the address, the Network is more easily used here.
		if protectedNet.Contains(net.ParseIP(addr.Network)) {
			if addressToKeep != nil {
				err = fmt.Errorf("duplicate protected addresses on vlan %d (%s)", vlanid, targIntf.Name)
				return
			}
			addressToKeep = &ipv4List[i]
		}
	}
	if addressToKeep == nil {
		err = fmt.Errorf("missing address in the 172.16.0.0/12 subnet on vlan %d (%s)", vlanid, targIntf.Name)
		return
	}

	// 	Find and remove all IPv4 routes where static=true and where gateway-status ends with <TargIntf>
	if routeList, err = conn.GetIPv4Routes([]string{"static"}); err != nil {
		return
	}
	for _, r := range routeList {
		if strings.HasSuffix(r.GatewayStatus, targIntf.Name) {
			fmt.Printf("Would remove IPv4 route %s\n", r.String())
			routesRemoved = append(routesRemoved, r)
		}
	}

	// If any removed routes were a /32 route, then see if that entry is in the PPP Secrets table
	if pppSecrets, err = conn.GetPPPSecrets(); err != nil {
		return
	}
	for _, s := range pppSecrets {
		if find32Route(routesRemoved, s.RemoteAddress) {
			s.Name = fmt.Sprintf("Static-%s", s.RemoteAddress)
			s.Comment = ""
			fmt.Printf("Would set PPP secret to %+v\n", s)
			ipsRemoved = append(ipsRemoved, routeros.IPv4Address{
				ID:      "",
				Address: s.RemoteAddress + "/32",
				Network: s.RemoteAddress,
			})
		}
	}

	// Remove any PPPoE servers for TargIntf
	if pppoeServers, e := conn.GetPPPoEServers(targIntf.Name); e != nil {
		err = e
		return
	} else {
		for i := range pppoeServers {
			fmt.Printf("Would remove PPPoE server %+v\n", pppoeServers[i])
		}
	}

	// Remove all IPv4 addresses except those in 172.16.0.0/12
	// Enable all remaining addresses (should just be the 172.16.0.0/12 one) if they were disabled
	for _, addr := range ipv4List {
		if addr.ID != addressToKeep.ID {
			fmt.Printf("Would remove IPv4 address %+v\n", addr)
			ipsRemoved = append(ipsRemoved, addr)
		} else if addr.Disabled {
			fmt.Printf("Would enable IPv4 address %+v\n", addr)
		}
	}

	// Find all queue tree entries where parent=<TargIntf>
	if queues, e := conn.GetQueueTree(targIntf.Name); e != nil {
		err = e
		return
	} else {
		for i := range queues {
			fmt.Printf("Would remove queue tree starting with %s\n", queues[i].Name)
		}
	}

	// Find all simply queue entries where target=<TargIntf>
	if simpQueues, e := conn.GetSimpleQueues(targIntf.Name); e != nil {
		err = e
		return
	} else {
		for _, q := range simpQueues {
			fmt.Printf("Would remove simply queue %+v\n", q)
		}
	}

	// Find DHCP server where interface=<TargIntf> and enable server if it was disabled
	if dhcpServers, e := conn.GetDhcp4ServerByIntf(targIntf.Name); e != nil {
		err = e
		return
	} else {
		if len(dhcpServers) > 1 {
			err = fmt.Errorf("more than one DHCPv4 server on vlan %d (%s)", vlanid, targIntf.Name)
			return
		}
		for _, s := range dhcpServers {
			if s.Disabled {
				fmt.Printf("Would enable DHCPv4 server on interface %s\n", s.Interface)
			}
		}
	}

	// Remove comment from TargIntf
	if len(targIntf.Comment) > 0 {
		fmt.Printf("Would remove comment '%s' from interface %s\n", targIntf.Comment, targIntf.Name)
	}

	// Rename TargIntf to `(^\w+-v\d{1,4}-)` + "avail"
	fmt.Printf("Would change interface name from %s to %s\n", targIntf.Name, nameParts[1]+"avail")
	return
}
