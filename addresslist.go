package gotik

import (
	"errors"
	"fmt"
	"net/netip"
	"strings"
	"time"
)

type AddressListAudit struct {
	Operation rune // A=Add, U=Update, D=Delete
	Address   string
	Comment   string
	ID        string
}

func (e *AddressListAudit) String() string {
	return fmt.Sprintf("%c: %s %s %s", e.Operation, e.ID, e.Address, e.Comment)
}

func (entry *AddressList) parse(props map[string]string) {
	for k, v := range props {
		switch k {
		case ".id":
			entry.ID = v
		case "list":
			entry.List = v
		case "address":
			entry.Address = v
		case "comment":
			entry.Comment = v
		case "dynamic":
			entry.Dynamic = parseBool(v)
		case "disabled":
			entry.Disabled = parseBool(v)
		case "creation-time":
			entry.CreationTime = parseTime(v)
		}
	}
	if entry.Dynamic {
		if t, found := props["timeout"]; found {
			entry.Timeout, _ = time.ParseDuration(t)
		}
	}
}

// GetIPv4AddressList returns a list of entries in an IPv4 address list
func (c *Client) GetIPv4AddressList(listname string) ([]AddressList, error) {
	return c.GetAddressList("ip", listname)
}

// GetIPv6AddressList returns a list of entries in an IPv4 address list
func (c *Client) GetIPv6AddressList(listname string) ([]AddressList, error) {
	return c.GetAddressList("ipv6", listname)
}

// GetIPv4AddressListEntry returns an entry in an IPv4 address list
func (c *Client) GetIPv4AddressListEntry(listname string, address string) (AddressList, error) {
	return c.GetAddressListEntry("ip", listname, address)
}

// GetIPv6AddressListEntry returns an entry in an IPv4 address list
func (c *Client) GetIPv6AddressListEntry(listname string, address string) (AddressList, error) {
	return c.GetAddressListEntry("ipv6", listname, address)
}

// GetAddressListEntry returns an entry in an address list.  The version must be "ipv4" or "ipv6"
func (c *Client) GetAddressListEntry(version string, listname string, address string) (AddressList, error) {
	var (
		err    error
		detail *Reply
		entry  AddressList
	)
	if version, err = deriveIPversion(version, address); err != nil {
		return entry, err
	}
	detail, err = c.Run("/"+version+"/firewall/address-list/print", "?=list="+listname, "?=address="+address)
	if err != nil {
		return entry, err
	}
	entry.parse(detail.Re[0].Map)
	return entry, nil
}

func (c *Client) GetAddressList(version string, listname string) ([]AddressList, error) {
	var (
		err    error
		detail *Reply
	)
	if version, err = deriveIPversion(version, ""); err != nil {
		return nil, err
	}
	detail, err = c.Run("/"+version+"/firewall/address-list/print", "?=list="+listname)
	if err != nil {
		return nil, err
	}
	entries := make([]AddressList, 0, 64)
	for _, re := range detail.Re {
		var entry AddressList
		entry.parse(re.Map)
		entries = append(entries, entry)
	}
	return entries, nil
}

// AuditIPv4AddressList audits an address list named by listname to ensure it has the addresses and comments as described
//
//		in the goodList map.  goodList must be a map which looks like:
//			map[string]string={
//				"62.4.108.20": "Location 1",
//				"10.0.0.0/8": "Management Network",
//			}
//		The existing list can be supplied in the list parameter, or nil can be used and the list will be retrieved.
//	   applyAudits must be true in order for changes to actually be applied to the router, else only proposed changes
//	   will be returned.  All proposed or executed changes are returned in the slice of AddressListAudit structs.
//	   If an empty goodList is supplied, the entire list will be removed from the router.
func (c *Client) AuditIPv4AddressList(listname string, list []AddressList, goodList map[string]string, applyAudits bool) ([]AddressListAudit, error) {
	// Adjust the goodList map to remove the /32 suffix on any entries. This matches what we will see in the list from the router.
	// Also force any FQDNs to lowercase
	targetList := make(map[string]string, len(goodList))
	for k, v := range goodList {
		if strings.HasSuffix(k, "/32") {
			targetList[k[0:len(k)-3]] = v
		} else {
			targetList[strings.ToLower(k)] = v
		}
	}
	return c.auditAddressList("ip", listname, list, targetList, applyAudits)
}

func (c *Client) AuditIPv6AddressList(listname string, list []AddressList, goodList map[string]string, applyAudits bool) ([]AddressListAudit, error) {
	// Force any FQDNs to lowercase
	targetList := make(map[string]string, len(goodList))
	for k, v := range goodList {
		targetList[strings.ToLower(k)] = v
	}
	return c.auditAddressList("ipv6", listname, list, targetList, applyAudits)
}

func (c *Client) auditAddressList(version string, listname string, list []AddressList, goodList map[string]string, applyAudits bool) ([]AddressListAudit, error) {
	var err error
	if len(listname) == 0 {
		return nil, fmt.Errorf("must supply a list name")
	}
	if list == nil {
		if list, err = c.GetAddressList(version, listname); err != nil {
			return nil, err
		}
	}

	// Allocate audit list to the maximum it might end up being (larger of goodList or list)
	maxLen := len(list)
	if len(goodList) > maxLen {
		maxLen = len(goodList)
	}
	audits := make([]AddressListAudit, 0, maxLen)

	// Execute Audit
	for _, e := range list {
		if e.Dynamic {
			continue
		}
		e.Address = strings.ToLower(e.Address) // in case it's a FQDN
		curAudit := AddressListAudit{Operation: 'N', ID: e.ID, Address: e.Address, Comment: e.Comment}
		comment, found := goodList[e.Address]
		if found && !e.Disabled { // if we found it, but it's disabled on the router, treat it as not found
			if e.Comment != comment {
				curAudit.Comment = comment
				curAudit.Operation = 'U'
			}
			delete(goodList, e.Address)
		} else {
			curAudit.Operation = 'D'
		}
		audits = append(audits, curAudit)
	}
	for k, v := range goodList {
		audits = append(audits, AddressListAudit{Operation: 'A', Address: k, Comment: v})
	}
	if applyAudits {
		for _, e := range audits {
			switch e.Operation {
			case 'A':
				_, err = c.Run(GenerateTikSentence("/"+version+"/firewall/address-list/add", "=", false,
					&AddressList{List: listname, Address: e.Address, Comment: e.Comment})...)
			case 'U':
				_, err = c.Run("/"+version+"/firewall/address-list/set", "=.id="+e.ID, "=comment="+e.Comment)
			case 'D':
				_, err = c.Run("/"+version+"/firewall/address-list/remove", "=.id="+e.ID)
			}
			if err != nil {
				break
			}
		}
	}
	return audits, err
}

// deriveIPversion will return "ip" or "ipv6" if the version isn't already one of those two values
// or if the supplied ip is an IPv4 or IPv6 address or subnet/prefix.
func deriveIPversion(version string, ip string) (string, error) {
	if version == "ip" || version == "ipv6" {
		return version, nil
	}
	if len(ip) == 0 {
		return "", errors.New("version must be ip or ipv6")
	}
	i := strings.LastIndexByte(ip, '/')
	if i < 0 {
		// not a prefix or subnet, so see if it's an IP address
		if a, err := netip.ParseAddr(ip); err == nil {
			if a.Is6() {
				return "ipv6", nil
			}
			return "ip", nil
		} else {
			return "", ErrIPorPrefixNotValid
		}
	}
	// Must be a subnet or prefix
	if a, err := netip.ParsePrefix(ip); err == nil {
		if a.Addr().Is6() {
			return "ipv6", nil
		} else {
			return "ip", nil
		}
	} else {
		return "", ErrIPorPrefixNotValid
	}
}

// AddIPv4AddressListEntry will add an entry to the named IP address list.  Address can be an IP, subnet, or DNS name.
// Timeout can be 0 if no timeout, or a sequence of positive decimal numbers, each with optional fraction and a
// unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// Comment can be an empty string if no comment is needed.  Disabled is a boolean but must be false if duration is not 0.
func (c *Client) AddIPv4AddressListEntry(listname string, address string, timeout time.Duration, comment string, disabled bool) (*Reply, error) {
	return c.AddAddressListEntry("ip", listname, address, timeout, comment, disabled)
}

// AddIPv6AddressListEntry will add an entry to the named IP address list.  Address can be an IP, subnet, or DNS name.
// Timeout can be 0 if no timeout, or a sequence of positive decimal numbers, each with optional fraction and a
// unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// Comment can be an empty string if no comment is needed.  Disabled is a boolean but must be false if duration is not 0.
func (c *Client) AddIPv6AddressListEntry(listname string, address string, timeout time.Duration, comment string, disabled bool) (*Reply, error) {
	return c.AddAddressListEntry("ipv6", listname, address, timeout, comment, disabled)
}

// AddAddressListEntry will add an entry to the named IP address list.  Address can be an IP, subnet, or DNS name.
// Timeout can be 0 if no timeout, or a sequence of positive decimal numbers, each with optional fraction and a
// unit suffix, such as "300ms", "-1.5h" or "2h45m". Valid time units are "ns", "us" (or "µs"), "ms", "s", "m", "h".
// Comment can be an empty string if no comment is needed.  Disabled is a boolean but must be false if duration is not 0.
func (c *Client) AddAddressListEntry(version string, listname string, address string, timeout time.Duration, comment string, disabled bool) (*Reply, error) {
	var (
		err   error
		reply *Reply
	)
	if version, err = deriveIPversion(version, address); err != nil {
		return nil, err
	}
	reply, err = c.Run(GenerateTikSentence("/"+version+"/firewall/address-list/add", "=", false,
		&AddressList{List: listname, Address: address, Timeout: timeout, Comment: comment, Disabled: disabled})...)
	if err != nil {
		if errors.Is(err, &ErrExists) {
			// Already exists, update
			var entry AddressList
			entry, err = c.GetAddressListEntry(version, listname, address)
			if err != nil {
				return nil, err
			}
			reply, err = c.Run("/"+version+"/firewall/address-list/set", "=.id="+entry.ID, "=comment="+comment,
				"=timeout="+timeout.Round(time.Second).String(), "=disabled="+fmt.Sprint(disabled))
			if err != nil {
				return reply, err
			}
		} else {
			return reply, err
		}
	}
	return reply, err
}

func (c *Client) RemoveAddressListEntry(version string, listname string, address string) error {
	var (
		err   error
		entry AddressList
	)
	if len(listname) == 0 {
		return errors.New("missing listname")
	}
	if len(address) == 0 {
		return errors.New("missing address")
	}
	if version, err = deriveIPversion(version, address); err != nil {
		return err
	}
	entry, err = c.GetAddressListEntry(version, listname, address)
	if err != nil {
		return err
	}
	_, err = c.Run("/"+version+"/firewall/address-list/remove", "=.id="+entry.ID)
	return err
}
