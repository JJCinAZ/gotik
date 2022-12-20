package gotik

import (
	"errors"
	"fmt"
	"strconv"
)

/*
Ethernet:
 (string) (len=3) ".id": (string) (len=2) "*2",
 (string) (len=10) "rx-128-255": (string) (len=6) "990446",
 (string) (len=8) "tx-bytes": (string) (len=11) "79040451056",
 (string) (len=11) "tx-1519-max": (string) (len=1) "0",
 (string) (len=8) "disabled": (string) (len=5) "false",
 (string) (len=12) "default-name": (string) (len=6) "ether1",
 (string) (len=3) "mtu": (string) (len=4) "1500",
 (string) (len=16) "auto-negotiation": (string) (len=4) "true",
 (string) (len=21) "tx-excessive-deferred": (string) (len=1) "0",
 (string) (len=11) "tx-too-long": (string) (len=1) "0",
 (string) (len=3) "arp": (string) (len=7) "enabled",
 (string) (len=11) "full-duplex": (string) (len=4) "true",
 (string) (len=16) "driver-tx-packet": (string) (len=8) "58308970",
 (string) (len=12) "rx-multicast": (string) (len=6) "364686",
 (string) (len=10) "tx-128-255": (string) (len=6) "980111",
 (string) (len=11) "mac-address": (string) (len=17) "E4:8D:8C:0F:D8:90",
 (string) (len=16) "driver-rx-packet": (string) (len=8) "21387556",
 (string) (len=19) "tx-single-collision": (string) (len=1) "0",
 (string) (len=7) "running": (string) (len=4) "true",
 (string) (len=5) "rx-64": (string) (len=7) "2671314",
 (string) (len=11) "rx-1519-max": (string) (len=1) "0",
 (string) (len=8) "rx-pause": (string) (len=1) "0",
 (string) (len=21) "tx-multiple-collision": (string) (len=1) "0",
 (string) (len=9) "advertise": (string) (len=59) "10M-half,10M-full,100M-half,100M-full,1000M-half,1000M-full",
 (string) (len=10) "tx-256-511": (string) (len=6) "854065",
 (string) (len=22) "tx-excessive-collision": (string) (len=1) "0",
 (string) (len=17) "tx-late-collision": (string) (len=1) "0",
 (string) (len=12) "loop-protect": (string) (len=7) "default",
 (string) (len=26) "loop-protect-send-interval": (string) (len=2) "5s",
 (string) (len=12) "rx-broadcast": (string) (len=6) "183960",
 (string) (len=14) "rx-align-error": (string) (len=1) "0",
 (string) (len=12) "rx-fcs-error": (string) (len=1) "0",
 (string) (len=9) "tx-65-127": (string) (len=7) "1724246",
 (string) (len=5) "slave": (string) (len=4) "true",
 (string) (len=11) "rx-512-1023": (string) (len=6) "527767",
 (string) (len=12) "rx-1024-1518": (string) (len=6) "814493",
 (string) (len=14) "driver-rx-byte": (string) (len=10) "3120293368",
 (string) (len=8) "rx-bytes": (string) (len=10) "3205843400",
 (string) (len=4) "name": (string) (len=6) "e1-lan",
 (string) (len=5) "l2mtu": (string) (len=4) "1578",
 (string) (len=15) "rx-flow-control": (string) (len=3) "off",
 (string) (len=12) "rx-too-short": (string) (len=1) "0",
 (string) (len=12) "tx-broadcast": (string) (len=5) "18712",
 (string) (len=8) "tx-pause": (string) (len=1) "0",
 (string) (len=11) "tx-deferred": (string) (len=1) "0",
 (string) (len=19) "loop-protect-status": (string) (len=3) "off",
 (string) (len=15) "tx-flow-control": (string) (len=3) "off",
 (string) (len=5) "speed": (string) (len=7) "100Mbps",
 (string) (len=10) "rx-256-511": (string) (len=6) "298650",
 (string) (len=11) "rx-too-long": (string) (len=1) "0",
 (string) (len=16) "orig-mac-address": (string) (len=17) "E4:8D:8C:0F:D8:90",
 (string) (len=5) "tx-64": (string) (len=6) "912573",
 (string) (len=11) "tx-512-1023": (string) (len=7) "1053093",
 (string) (len=12) "tx-1024-1518": (string) (len=8) "52784879",
 (string) (len=12) "tx-multicast": (string) (len=5) "57959",
 (string) (len=25) "loop-protect-disable-time": (string) (len=2) "5m",
 (string) (len=14) "driver-tx-byte": (string) (len=11) "78800125158",
 (string) (len=11) "tx-underrun": (string) (len=1) "0",
 (string) (len=11) "arp-timeout": (string) (len=4) "auto",
 (string) (len=9) "rx-65-127": (string) (len=8) "16084884",
 (string) (len=9) "bandwidth": (string) (len=19) "unlimited/unlimited",
 (string) (len=6) "switch": (string) (len=7) "switch1",
 (string) (len=11) "rx-fragment": (string) (len=1) "0",
 (string) (len=11) "rx-overflow": (string) (len=1) "0",
 (string) (len=12) "tx-collision": (string) (len=1) "0"

Bridge:
 (string) (len=3) ".id": (string) (len=2) "*6",
 (string) (len=3) "mtu": (string) (len=4) "1500",
 (string) (len=5) "l2mtu": (string) (len=5) "65535",
 (string) (len=11) "mac-address": (string) (len=17) "02:43:D4:21:14:00",
 (string) (len=7) "running": (string) (len=4) "true",
 (string) (len=9) "admin-mac": (string) (len=17) "02:43:D4:21:14:00",
 (string) (len=13) "dhcp-snooping": (string) (len=5) "false",
 (string) (len=10) "actual-mtu": (string) (len=4) "1500",
 (string) (len=11) "arp-timeout": (string) (len=4) "auto",
 (string) (len=13) "protocol-mode": (string) (len=4) "none",
 (string) (len=13) "igmp-snooping": (string) (len=5) "false",
 (string) (len=8) "auto-mac": (string) (len=5) "false",
 (string) (len=4) "name": (string) (len=3) "lo0",
 (string) (len=11) "ageing-time": (string) (len=2) "5m",
 (string) (len=14) "vlan-filtering": (string) (len=5) "false",
 (string) (len=8) "disabled": (string) (len=5) "false",
 (string) (len=3) "arp": (string) (len=7) "enabled",
 (string) (len=12) "fast-forward": (string) (len=5) "false"

Generic Interface: (/interface print)
 (string) (len=9) "rx-packet": (string) (len=9) "443637812",
 (string) (len=7) "rx-drop": (string) (len=1) "0",
 (string) (len=12) "fp-rx-packet": (string) (len=1) "0",
 (string) (len=12) "fp-tx-packet": (string) (len=1) "0",
 (string) (len=3) "mtu": (string) (len=4) "1500",
 (string) (len=11) "mac-address": (string) (len=17) "4C:5E:0C:0F:FA:2D",
 (string) (len=10) "link-downs": (string) (len=1) "0",
 (string) (len=7) "rx-byte": (string) (len=12) "174918486575",
 (string) (len=8) "tx-error": (string) (len=1) "0",
 (string) (len=10) "fp-tx-byte": (string) (len=1) "0",
 (string) (len=4) "name": (string) (len=16) "e1-v195-Nextrio2",
 (string) (len=4) "type": (string) (len=4) "vlan",
 (string) (len=10) "actual-mtu": (string) (len=4) "1500",
 (string) (len=7) "tx-byte": (string) (len=12) "608155461130",
 (string) (len=13) "tx-queue-drop": (string) (len=1) "0",
 (string) (len=8) "rx-error": (string) (len=1) "0",
 (string) (len=10) "fp-rx-byte": (string) (len=1) "0",
 (string) (len=3) ".id": (string) (len=2) "*8",
 (string) (len=17) "last-link-up-time": (string) (len=20) "jul/22/2019 19:53:52",
 (string) (len=9) "tx-packet": (string) (len=9) "514759968",
 (string) (len=7) "tx-drop": (string) (len=1) "0",
 (string) (len=7) "running": (string) (len=4) "true",
 (string) (len=8) "disabled": (string) (len=5) "false",
 (string) (len=5) "l2mtu": (string) (len=4) "1576"
}
2019/09/18 15:16:48 (map[string]string) (len=24) {
 (string) (len=7) "tx-byte": (string) (len=8) "37459260",
 (string) (len=9) "tx-packet": (string) (len=6) "249730",
 (string) (len=7) "rx-drop": (string) (len=1) "0",
 (string) (len=10) "fp-tx-byte": (string) (len=1) "0",
 (string) (len=7) "running": (string) (len=4) "true",
 (string) (len=5) "l2mtu": (string) (len=5) "65535",
 (string) (len=10) "link-downs": (string) (len=1) "0",
 (string) (len=11) "mac-address": (string) (len=17) "02:43:D4:21:14:00",
 (string) (len=13) "tx-queue-drop": (string) (len=1) "0",
 (string) (len=8) "rx-error": (string) (len=1) "0",
 (string) (len=10) "fp-rx-byte": (string) (len=1) "0",
 (string) (len=12) "fp-rx-packet": (string) (len=1) "0",
 (string) (len=4) "type": (string) (len=6) "bridge",
 (string) (len=10) "actual-mtu": (string) (len=4) "1500",
 (string) (len=8) "tx-error": (string) (len=1) "0",
 (string) (len=17) "last-link-up-time": (string) (len=20) "jul/22/2019 19:53:42",
 (string) (len=9) "rx-packet": (string) (len=1) "0",
 (string) (len=3) "mtu": (string) (len=4) "1500",
 (string) (len=7) "rx-byte": (string) (len=1) "0",
 (string) (len=7) "tx-drop": (string) (len=1) "0",
 (string) (len=12) "fp-tx-packet": (string) (len=1) "0",
 (string) (len=8) "disabled": (string) (len=5) "false",
 (string) (len=3) ".id": (string) (len=2) "*6",
 (string) (len=4) "name": (string) (len=3) "lo0"
}
*/

/*
reply, err = routerConn.Run("/interface/ethernet/monitor", "=.id=*1,*2,*3,*4", "=once=yes")
Ethernet Monitor returns:
 {`name` `ge03`}
 {`status` `link-ok`}
 {`auto-negotiation` `done`}
 {`rate` `100Mbps`}
 {`full-duplex` `true`}
 {`tx-flow-control` `false`}
 {`rx-flow-control` `false`}
 {`advertising` `10M-half,10M-full,100M-half,100M-full,1000M-full`}
 {`link-partner-advertising` `10M-half,10M-full,100M-half,100M-full`}
*/

func parseInterface(props map[string]string, intfType string) Interface {
	var found bool
	intf := Interface{
		ID:        props[".id"],
		Type:      props["type"],
		Name:      props["name"],
		Comment:   props["comment"],
		Mac:       props["mac-address"],
		Interface: props["interface"],
		Disabled:  parseBool(props["disabled"]),
		Dynamic:   parseBool(props["dynamic"]),
		Running:   parseBool(props["running"]),
		Arp:       props["arp"],
		AutoNeg:   parseBool(props["auto-negotiation"]),
		MTU:       parseInt(props["mtu"]),
		ActualMTU: parseInt(props["actual-mtu"]),
	}
	if intf.Type, found = props["type"]; !found {
		intf.Type = intfType
	}
	switch intf.Type {
	case InterfaceTypeVlan:
		x, _ := strconv.ParseInt(props["vlan-id"], 10, 64)
		intf.VLAN = int(x)
	case InterfaceTypeEthernet:
		if intf.AutoNeg == false {
			// The speed prop isn't valid unless we are doing manual negotiation
			intf.Speed = props["speed"]
		}
		intf.DefaultName = props["default-name"]
		intf.OriginalMac = props["orig-mac-address"]
		intf.Slave = parseBool(props["slave"])
	case InterfaceTypeBridge:
		intf.AdminMac = props["admin-mac"]
		intf.AutoMac = parseBool(props["auto-mac"])
		intf.ProtocolMode = props["protocol-mode"]
		intf.AgingTime = props["ageing-time"]
		intf.VlanFiltering = parseBool(props["vlan-filtering"])
		intf.FastForward = parseBool(props["fast-forward"])
	case InterfaceTypeGre:
		intf.LocalAddress = props["local-address"]
		intf.RemoteAddress = props["remote-address"]
		intf.KeepAlive = props["keepalive"]
		intf.IPSecSecret = props["ipsec-secret"]
	}
	return intf
}

// SetInterfaceComment changes the comment on an interface.  This can be used for any interface type
func (c *Client) SetInterfaceComment(id string, comment string) error {
	_, err := c.Run("/interface/set", "=.id="+id, "=comment="+comment)
	return err
}

func (c *Client) SetInterfaceName(id string, newName string) error {
	_, err := c.Run("/interface/set", "=.id="+id, "=name="+newName)
	return err
}

func (c *Client) EnableInterface(id string) error {
	_, err := c.Run("/interface/enable", "=.id="+id)
	return err
}

func (c *Client) DisableInterface(id string) error {
	_, err := c.Run("/interface/disable", "=.id="+id)
	return err
}

// GetVLANInterface returns a single VLAN interface, or an error if the VLAN ID is not unique on the router
func (c *Client) GetVLANInterface(vlan int) (intf Interface, err error) {
	detail, err := c.Run("/interface/vlan/print", fmt.Sprintf("?=vlan-id=%d", vlan))
	if err != nil {
		return
	}

	// only return the id if there is one result
	switch len(detail.Re) {
	case 0:
		err = ErrNotFound
	case 1:
		intf = parseInterface(detail.Re[0].Map, InterfaceTypeVlan)
	default:
		err = errors.New("more than one matching interface - only expecting one")
	}
	return
}

// AddVLANInterface adds a new VLAN interface.
// The parameter intf.Type must be InterfaceTypeVlan
func (c *Client) AddVLANInterface(intf Interface) (string, error) {
	if intf.Type != InterfaceTypeVlan || len(intf.Name) == 0 || len(intf.Interface) == 0 {
		return "", fmt.Errorf("invalid interface supplied")
	}
	if intf.VLAN < 1 || intf.VLAN > 4095 {
		return "", fmt.Errorf("vlan ID out of range (1...4095)")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/interface/vlan/add")
	parts = append(parts, fmt.Sprintf("=name=%s", intf.Name))
	parts = append(parts, fmt.Sprintf("=disabled=%t", intf.Disabled))
	parts = append(parts, fmt.Sprintf("=vlan-id=%d", intf.VLAN))
	parts = append(parts, fmt.Sprintf("=interface=%s", intf.Interface))
	if len(intf.Arp) > 0 {
		parts = append(parts, fmt.Sprintf("=arp=%s", intf.Arp))
	}
	if len(intf.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", intf.Comment))
	}
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

// GetVLANInterfaceOnBase returns a single VLAN on a base interface or an error if the VLAN is not found
func (c *Client) GetVLANInterfaceOnBase(baseIntf string, vlan int) (intf Interface, err error) {
	detail, err := c.Run("/interface/vlan/print",
		fmt.Sprintf("?=vlan-id=%d", vlan),
		fmt.Sprintf("?=interface="+baseIntf))
	if err != nil {
		return
	}

	// only return the id if there is one result
	switch len(detail.Re) {
	default:
		err = ErrNotFound
	case 1:
		intf = parseInterface(detail.Re[0].Map, InterfaceTypeVlan)
	}
	return
}

// GetVlanInterfaces returns a list of all VLAN interfaces on a particular base interface or all interfaces if baseIntf is blank
func (c *Client) GetVlanInterfaces(baseIntf string) ([]Interface, error) {
	var (
		err    error
		detail *Reply
	)
	if len(baseIntf) > 0 {
		detail, err = c.Run("/interface/vlan/print", "?=interface="+baseIntf)
	} else {
		detail, err = c.Run("/interface/vlan/print")
	}
	if err != nil {
		return nil, err
	}
	interfaces := make([]Interface, 0, 256)
	for _, re := range detail.Re {
		interfaces = append(interfaces, parseInterface(re.Map, InterfaceTypeVlan))
	}
	return interfaces, nil
}

// GetEthInterfaces returns a list of all Ethernet interfaces
func (c *Client) GetEthInterfaces() ([]Interface, error) {
	detail, err := c.Run("/interface/ethernet/print")
	if err != nil {
		return nil, err
	}
	interfaces := make([]Interface, 0, 32)
	for _, re := range detail.Re {
		interfaces = append(interfaces, parseInterface(re.Map, InterfaceTypeEthernet))
	}
	return interfaces, nil
}

// GetBridgeInterfaces returns a list of all Bridge interfaces
func (c *Client) GetBridgeInterfaces() ([]Interface, error) {
	detail, err := c.Run("/interface/bridge/print")
	if err != nil {
		return nil, err
	}
	interfaces := make([]Interface, 0, 32)
	for _, re := range detail.Re {
		interfaces = append(interfaces, parseInterface(re.Map, InterfaceTypeBridge))
	}
	return interfaces, nil
}

// GetInterfacesOfTypes returns a list of all Bridge interfaces
func (c *Client) GetInterfacesOfTypes(types ...string) ([]Interface, error) {
	detail, err := c.Run("/interface/print")
	if err != nil {
		return nil, err
	}
	interfaces := make([]Interface, 0, 64)
	for _, re := range detail.Re {
		for _, t := range types {
			if t == re.Map["type"] {
				interfaces = append(interfaces, parseInterface(re.Map, t))
				break
			}
		}
	}
	return interfaces, nil
}
