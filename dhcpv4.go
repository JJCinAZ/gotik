package gotik

import (
	"fmt"
)

/*
address -- Network address
boot-file-name -- Boot file name
caps-manager --
comment -- Short description of the item
copy-from -- Item number
dhcp-option --
dhcp-option-set --
dns-server -- DNS server address
domain -- Domain name
gateway -- Default gateway
netmask -- Network mask to give out
next-server -- IP address of next server to use in bootstrap
ntp-server -- NTP server
wins-server -- WINS server
*/
func parseDhcp4Network(props map[string]string) DHCP4Network {
	entry := DHCP4Network{
		ID:         props[".id"],
		Address:    props["address"],
		Netmask:    props["netmask"],
		Gateway:    props["gateway"],
		Domain:     props["domain"],
		Comment:    props["comment"],
		DNSServers: props["dns-server"],
		NTPServers: props["ntp-server"],
		Options:    props["dhcp-options"],
	}
	return entry
}

/*
add-arp -- Defines whether to add dynamic ARP entry
address-pool -- IP address pool
always-broadcast -- Send all replies as broadcast
authoritative -- This is the only DHCP server for the network
bootp-lease-time --
bootp-support -- Support for BOOTP clients
conflict-detection --
copy-from -- Item number
delay-threshold -- If secs field in DHCP packet is smaller than delay-threshold, then this packet is ignored
disabled -- Defines whether item is ignored or used
insert-queue-before --
interface -- Interface
lease-script --
lease-time -- Lease time
name -- Name of DHCP server
relay -- DHCP relay address
src-address -- Source address
use-framed-as-classless --
use-radius -- Use RADIUS server for authentication
*/
func parseDhcp4Server(props map[string]string) DHCPv4Server {
	entry := DHCPv4Server{
		ID:            props[".id"],
		Name:          props["name"],
		LeaseTime:     props["lease-time"],
		Pool:          props["address-pool"],
		Interface:     props["interface"],
		Authoritative: props["authoritative"],
		Disabled:      parseBool(props["disabled"]),
		AddArp:        parseBool(props["add-arp"]),
	}
	return entry
}

func (c *Client) dhcp4ServerPrint(parms ...string) ([]DHCPv4Server, error) {
	entries := make([]DHCPv4Server, 0, 8)
	detail, err := c.RunCmd("/ip/dhcp-server/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseDhcp4Server(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// Returns a single DHCP Server by name
func (c *Client) GetDhcp4ServerByIntf(intf string) ([]DHCPv4Server, error) {
	return c.dhcp4ServerPrint("?=interface=" + intf)
}

// Returns a single DHCP Server by name
func (c *Client) GetDhcp4ServerByName(name string) ([]DHCPv4Server, error) {
	return c.dhcp4ServerPrint("?=name=" + name)
}

// Returns a list of all DHCP Servers
func (c *Client) GetDhcpv4Servers() ([]DHCPv4Server, error) {
	return c.dhcp4ServerPrint()
}

// Add a new DHCPv4 Server
func (c *Client) AddDhcpv4Server(s DHCPv4Server) (string, error) {
	if len(s.Name) == 0 || len(s.Interface) == 0 || len(s.Pool) == 0 || len(s.LeaseTime) == 0 {
		return "", fmt.Errorf("invalid DHCPv4 server supplied")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/ip/dhcp-server/add")
	parts = append(parts, fmt.Sprintf("=name=%s", s.Name))
	parts = append(parts, fmt.Sprintf("=interface=%s", s.Interface))
	parts = append(parts, fmt.Sprintf("=address-pool=%s", s.Pool))
	parts = append(parts, fmt.Sprintf("=lease-time=%s", s.LeaseTime))
	parts = append(parts, fmt.Sprintf("=disabled=%t", s.Disabled))
	parts = append(parts, fmt.Sprintf("=add-arp=%t", s.AddArp))
	if len(s.Authoritative) > 0 {
		parts = append(parts, fmt.Sprintf("=authoritative=%s", s.Authoritative))
	}
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (c *Client) SetDhcpv4ServerDisable(id string, disabled bool) error {
	d := "=disabled=true"
	if !disabled {
		d = "=disabled=false"
	}
	_, err := c.Run("/ip/dhcp-server/set", "=.id="+id, d)
	if err != nil {
		return err
	}
	return nil
}

func (c *Client) dhcp4NetworkPrint(parms ...string) ([]DHCP4Network, error) {
	entries := make([]DHCP4Network, 0, 8)
	detail, err := c.RunCmd("/ip/dhcp-server/network/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseDhcp4Network(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// Returns a list of all DHCP Networks
func (c *Client) GetDhcpv4Networks() ([]DHCP4Network, error) {
	return c.dhcp4NetworkPrint()
}

// Add a new DHCPv4 Network
func (c *Client) AddDhcpv4Network(s DHCP4Network) (string, error) {
	if len(s.Address) == 0 || len(s.Netmask) == 0 {
		return "", fmt.Errorf("invalid DHCPv4 network")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/ip/dhcp-server/network/add")
	parts = append(parts, fmt.Sprintf("=address=%s", s.Address))
	parts = append(parts, fmt.Sprintf("=netmask=%s", s.Netmask))
	if len(s.Gateway) > 0 {
		parts = append(parts, fmt.Sprintf("=gateway=%s", s.Gateway))
	}
	if len(s.Domain) > 0 {
		parts = append(parts, fmt.Sprintf("=domain=%s", s.Domain))
	}
	if len(s.DNSServers) > 0 {
		parts = append(parts, fmt.Sprintf("=dns-server=%s", s.DNSServers))
	}
	if len(s.NTPServers) > 0 {
		parts = append(parts, fmt.Sprintf("=ntp-server=%s", s.NTPServers))
	}
	if len(s.Options) > 0 {
		parts = append(parts, fmt.Sprintf("=dhcp-option=%s", s.Options))
	}
	if len(s.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", s.Comment))
	}
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}
