package gotik

import (
	"errors"
	"fmt"
	"net"
)

func parsev4Addr(props map[string]string) IPv4Address {
	entry := IPv4Address{
		ID:              props[".id"],
		Address:         props["address"],
		Network:         props["network"],
		Interface:       props["interface"],
		ActualInterface: props["actual-interface"],
		Comment:         props["comment"],
		Invalid:         parseBool(props["invalid"]),
		Dynamic:         parseBool(props["dynamic"]),
		Disabled:        parseBool(props["disabled"]),
	}
	if len(entry.Address) > 3 {
		// special decoding based on subnet type
		switch entry.Address[len(entry.Address)-3:] {
		case "/32":
			// change subnet to /31 if the address and network are different - indicating a /31 on a tik
			if entry.Address != entry.Network {
				entry.Address = entry.Address[:len(entry.Address)-2] + "31"
			}
		}
	}
	return entry
}

func (c *Client) ipAddressPrint(parms ...string) ([]IPv4Address, error) {
	entries := make([]IPv4Address, 0, 8)
	detail, err := c.RunCmd("/ip/address/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parsev4Addr(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// GetInterfaceIPv4Table returns a list of all IPv4 addresses on a particular interface
func (c *Client) GetInterfaceIPv4Table(baseIntf string) ([]IPv4Address, error) {
	return c.ipAddressPrint("?=interface=" + baseIntf)
}

// GetIPv4Table returns a list of all IPv4 addresses on the router
func (c *Client) GetIPv4Table() ([]IPv4Address, error) {
	return c.ipAddressPrint()
}

// AddIPv4Address adds a new IPv4 Address
func (c *Client) AddIPv4Address(addr IPv4Address) (string, error) {
	if len(addr.Address) == 0 || len(addr.Interface) == 0 {
		return "", fmt.Errorf("invalid IPv4 address supplied")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/ip/address/add")
	parts = append(parts, fmt.Sprintf("=address=%s", addr.Address))
	parts = append(parts, fmt.Sprintf("=interface=%s", addr.Interface))
	if len(addr.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", addr.Comment))
	}
	parts = append(parts, fmt.Sprintf("=disabled=%t", addr.Disabled))
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

// ModifyIPv4Address changes the specified property of an address
func (c *Client) ModifyIPv4Address(id string, action string) error {
	switch action {
	case "enable":
		_, err := c.Run("/ip/address/set", "=disabled=no", "=.id="+id)
		if err != nil {
			return err
		}
	case "disable":
		_, err := c.Run("/ip/address/set", "=disabled=yes", "=.id="+id)
		if err != nil {
			return err
		}
	case "remove":
		_, err := c.Run("/ip/address/remove", "=.id="+id)
		if err != nil {
			return err
		}
	default:
		return errors.New("address modification action invalid")
	}

	// return nil if all good
	return nil
}

// GetCustomerIPv4Subnets returns a slice of all addresses associated with a specified vlan
// looks into addresses on the interface and any static routes to them
func (c *Client) GetCustomerIPv4Subnets(vlan int) ([]IPv4Address, error) {
	var routes []IPv4Route

	// get vlan interface info, returns error if vlan is not unique on router
	interfaceInfo, err := c.GetVLANInterface(vlan)
	if err != nil {
		return nil, err
	}

	// find addresses for target vlan interface
	addresses, err := c.GetInterfaceIPv4Table(interfaceInfo.Name)
	if err != nil {
		return nil, err
	}

	// handle /32 addresses setup in static routes
	staticRoutes, err := c.GetIPv4Routes([]string{"static"})
	if err != nil {
		return nil, err
	}

	// find routes matching interface addresses
	for i := range addresses {
		/*
			If Address looks something like: "10.20.3.5/30" then Network looks like "10.20.3.4"
			We need the Network to look like "10.20.3.4/30" for our call to FindMatchingRoutes
		*/
		_, addressNet, err := net.ParseCIDR(addresses[i].Address)
		if err != nil {
			return nil, err
		}
		addresses[i].Network = addressNet.String()

		// look for static routes where the gateway falls into the address subnet
		matchedRoutes, err := c.FindMatchingIPv4Routes(staticRoutes, addressNet)
		if err != nil {
			return nil, err
		}

		// add matching routes to the return
		if len(matchedRoutes) > 0 {
			routes = append(routes, matchedRoutes...)
		}
	}

	// pull destination addresses from matched routes
	for i := range routes {
		address, addressNet, err := net.ParseCIDR(routes[i].DstAddress)
		if err != nil {
			return nil, err
		}
		addresses = append(addresses, IPv4Address{Address: address.String(), Network: addressNet.String()})
	}

	// TODO handle PPPoE interfaces

	// return addresses if all good
	return addresses, nil
}
