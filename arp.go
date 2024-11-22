package gotik

import (
	"fmt"
	"net"
	"regexp"
	"strings"
)

var regexpMacAddressSeparated = regexp.MustCompile(`^([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})[:\-]?([[:xdigit:]]{2})$`)

// NormalizeMac will take a MAC address in one of the formats: ############, ##-##-##-##-##-## or ##:##:##:##:##:##
// and will return it normalized to upper case with the specified separator.
func NormalizeMac(s string, separator string) (string, error) {
	if m := regexpMacAddressSeparated.FindStringSubmatch(strings.ToUpper(strings.TrimSpace(s))); m != nil && len(m) == 7 {
		return strings.Join(m[1:], separator), nil
	}
	return "", fmt.Errorf("invalid format for MAC address '%-1.20s'", s)
}

func parseArp(props map[string]string) ArpEntry {
	entry := ArpEntry{
		ID:        props[".id"],
		Address:   props["address"],
		Mac:       props["mac-address"],
		Interface: props["interface"],
		Comment:   props["comment"],
		Disabled:  parseBool(props["disabled"]),
		Dynamic:   parseBool(props["dynamic"]),
		Complete:  parseBool(props["complete"]),
		DHCP:      parseBool(props["DHCP"]),
	}
	return entry
}

func (c *Client) ipArpPrint(parms ...string) ([]ArpEntry, error) {
	entries := make([]ArpEntry, 0, 256)
	detail, err := c.RunCmd("/ip/arp/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseArp(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// GetInterfaceArpTable returns a list of all ARP entries on a particular interface
func (c *Client) GetInterfaceArpTable(baseIntf string) ([]ArpEntry, error) {
	return c.ipArpPrint("?=interface=" + baseIntf)
}

// GetArpTable returns a list of all ARP entries
func (c *Client) GetArpTable() ([]ArpEntry, error) {
	return c.ipArpPrint()
}

// ArpLookupByIP returns any entry for a particular IP address
func (c *Client) ArpLookupByIP(ipv4 string) (entry ArpEntry, err error) {
	if len(net.ParseIP(ipv4).To4()) != net.IPv4len {
		err = fmt.Errorf("Invalid IPv4 address")
		return
	}
	detail, err := c.Run("/ip/arp/print", "?=disabled=false", "?=address="+ipv4)
	if err != nil {
		return
	}
	if len(detail.Re) >= 1 {
		entry = parseArp(detail.Re[0].Map)
		return
	}
	err = ErrNotFound
	return
}

// ArpLookupByMAC returns any entry for a particular MAC address
func (c *Client) ArpLookupByMAC(mac string) (entry ArpEntry, err error) {
	var target string
	target, err = NormalizeMac(mac, ":")
	detail, err := c.Run("/ip/arp/print", "?=disabled=false", "?=mac-address="+target)
	if err != nil {
		return
	}
	if len(detail.Re) >= 1 {
		entry = parseArp(detail.Re[0].Map)
		return
	}
	err = ErrNotFound
	return
}
