package gotik

import (
	"fmt"
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

// Returns a list of entries in an IPv4 address list
func (c *Client) GetIPv4AddressList(listname string) ([]AddressList, error) {
	return c.getAddressList("/ip", listname)
}

// Returns a list of entries in an IPv4 address list
func (c *Client) GetIPv6AddressList(listname string) ([]AddressList, error) {
	return c.getAddressList("/ipv6", listname)
}

func (c *Client) getAddressList(cmdprefix string, listname string) ([]AddressList, error) {
	detail, err := c.Run(cmdprefix+"/firewall/address-list/print", "?=list="+listname)
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
//	in the goodList map.  goodList must be a map which looks like:
//		map[string]string={
//			"62.4.108.20": "Location 1",
//			"10.0.0.0/8": "Management Network",
//		}
//	The existing list can be supplied in the list parameter, or nil can be used and the list will be retrieved.
//    applyAudits must be true in order for changes to actually be applied to the router, else only proposed changes
//    will be returned.  All proposed or executed changes are returned in the slice of AddressListAudit structs.
//    If an empty goodList is supplied, the entire list will be removed from the router.
//
func (c *Client) AuditIPv4AddressList(listname string, list []AddressList, goodList map[string]string, applyAudits bool) ([]AddressListAudit, error) {
	return c.auditAddressList("/ip", listname, list, goodList, applyAudits)
}

func (c *Client) AuditIPv6AddressList(listname string, list []AddressList, goodList map[string]string, applyAudits bool) ([]AddressListAudit, error) {
	return c.auditAddressList("/ipv6", listname, list, goodList, applyAudits)
}

func (c *Client) auditAddressList(cmdprefix string, listname string, list []AddressList, goodList map[string]string, applyAudits bool) ([]AddressListAudit, error) {
	var err error
	if len(listname) == 0 {
		return nil, fmt.Errorf("must supply a list name")
	}
	if list == nil {
		if list, err = c.getAddressList(cmdprefix, listname); err != nil {
			return nil, err
		}
	}
	maxLen := len(list)
	if len(goodList) > maxLen {
		maxLen = len(goodList)
	}
	audits := make([]AddressListAudit, 0, maxLen)
	for _, e := range list {
		curAudit := AddressListAudit{Operation: 'N', ID: e.ID, Address: e.Address, Comment: e.Comment}
		if comment, found := goodList[e.Address]; found {
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
				_, err = c.Run(GenerateTikSentence(cmdprefix+"/firewall/address-list/add", "=", false,
					&AddressList{List: listname, Address: e.Address, Comment: e.Comment})...)
			case 'U':
				_, err = c.Run(cmdprefix+"/firewall/address-list/set", "=.id="+e.ID, "=comment="+e.Comment)
			case 'D':
				_, err = c.Run(cmdprefix+"/firewall/address-list/remove", "=.id="+e.ID)
			}
			if err != nil {
				break
			}
		}
	}
	return audits, err
}
