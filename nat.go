package gotik

import (
	"fmt"
	"strings"
)

func (c *Client) GetIPv4Nat(chain string) ([]IPv4NatRule, error) {
	a := []string{"/ip/firewall/nat/print"}
	if len(chain) > 0 {
		a = append(a, "?=chain="+chain)
	}
	detail, err := c.Run(a...)
	if err != nil {
		return nil, err
	}
	entries := make([]IPv4NatRule, 0, 64)
	for _, re := range detail.Re {
		var entry IPv4NatRule
		parseTikObject(re.Map, &entry)
		entries = append(entries, entry)
	}
	return entries, nil
}

func (c *Client) RemoveIPv4NatRule(id string) error {
	if len(id) == 0 {
		return ErrMissingId
	}
	_, err := c.Run("/ip/firewall/nat/remove", "=.id="+id)
	return err
}

func (r *IPv4NatRule) String() string {
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
	if r.SrcPort > 0 {
		a = append(a, fmt.Sprintf("src-port=%d", r.SrcPort))
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
	if r.DstPort > 0 {
		a = append(a, fmt.Sprintf("dst-port=%d", r.DstPort))
	}
	if r.Action == "jump" {
		a = append(a, "action=jump", fmt.Sprintf("jump-target=%s", r.JumpTarget))
	} else if r.Action == "return" {
		a = append(a, "action=return")
	} else {
		a = append(a, fmt.Sprintf("action=%s", r.Action), fmt.Sprintf("to-addresses=%s", r.ToAddresses))
		if r.ToPorts != 0 {
			a = append(a, fmt.Sprintf("to-ports=%d", r.ToPorts))
		}
	}
	if len(r.Comment) > 0 {
		a = append(a, fmt.Sprintf("comment=%s", r.Comment))
	}
	return strings.Join(a, " ")
}
