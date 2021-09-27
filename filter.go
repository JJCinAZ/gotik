package gotik

import (
	"fmt"
	"strings"
	"time"
)

func (c *Client) GetIPv4Filters(chain string) ([]IPv4FilterRule, error) {
	a := []string{"/ip/firewall/filter/print"}
	if len(chain) > 0 {
		a = append(a, "?=chain="+chain)
	}
	detail, err := c.Run(a...)
	if err != nil {
		return nil, err
	}
	entries := make([]IPv4FilterRule, 0, 64)
	for _, re := range detail.Re {
		var entry IPv4FilterRule
		parseTikObject(re.Map, &entry)
		// Correct for bug in certain ROS versions where "action=accept" is returned as an empty action
		if len(entry.Action) == 0 {
			entry.Action = "accept"
		}
		entries = append(entries, entry)
	}
	return entries, nil
}

func (c *Client) RemoveIPv4FilterRule(id string) error {
	if len(id) == 0 {
		return ErrMissingId
	}
	_, err := c.Run("/ip/firewall/filter/remove", "=.id="+id)
	return err
}

func (r *IPv4FilterRule) String() string {
	a := make([]string, 0, 24)
	f := []rune{' ', ' ', ' '}
	if r.Disabled {
		f[0] = 'X'
	}
	if r.Invalid {
		f[1] = 'I'
	}
	if r.Dynamic {
		f[2] = 'D'
	}
	a = append(a, string(f), fmt.Sprintf("chain=%s", r.Chain))
	if len(r.Protocol) > 0 {
		a = append(a, fmt.Sprintf("protocol=%s", r.Protocol))
	}
	if len(r.SrcAddressList) > 0 {
		a = append(a, fmt.Sprintf("src-address-list=%s", r.SrcAddressList))
	}
	if len(r.SrcAddress) > 0 {
		a = append(a, fmt.Sprintf("src-address=%s", r.SrcAddress))
	}
	if len(r.SrcPort) > 0 {
		a = append(a, fmt.Sprintf("src-port=%s", r.SrcPort))
	}
	if len(r.InInterface) > 0 {
		a = append(a, fmt.Sprintf("in-interface=%s", r.InInterface))
	}
	if len(r.OutInterface) > 0 {
		a = append(a, fmt.Sprintf("out-interface=%s", r.OutInterface))
	}
	if len(r.DstAddressList) > 0 {
		a = append(a, fmt.Sprintf("dst-address-list=%s", r.DstAddressList))
	}
	if len(r.DstAddress) > 0 {
		a = append(a, fmt.Sprintf("dst-address=%s", r.DstAddress))
	}
	if len(r.DstPort) > 0 {
		a = append(a, fmt.Sprintf("dst-port=%s", r.DstPort))
	}
	if r.Action == "jump" {
		a = append(a, fmt.Sprintf("action=%s", r.Action), fmt.Sprintf("jump-target=%s", r.JumpTarget))
	} else if r.Action == "add-dst-to-address-list" || r.Action == "add-dsrc-to-address-list" {
		a = append(a, fmt.Sprintf("action=%s", r.Action), fmt.Sprintf("address-list=%s", r.AddressList))
		if r.AddressListTimeout != 0 {
			a = append(a, fmt.Sprintf("address-list-timeout=%d seconds", r.AddressListTimeout/time.Second))
		}
	} else {
		a = append(a, fmt.Sprintf("action=%s", r.Action))
	}
	if len(r.Comment) > 0 {
		a = append(a, fmt.Sprintf("comment=%s", r.Comment))
	}
	return strings.Join(a, " ")
}
