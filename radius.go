package gotik

import (
	"fmt"
	"strings"
	"time"
)

func (e *RadiusServer) String() string {
	flags := ' '
	if e.Disabled {
		flags = 'X'
	}
	return fmt.Sprintf("%c %s: %s %s %s", flags, e.Service, e.Address, e.Secret, e.Comment)
}

func parseRadius(props map[string]string) RadiusServer {
	entry := RadiusServer{
		ID:                 props[".id"],
		AccountingBackup:   parseBool(props["accounting-backup"]),
		AccountingPort:     parseInt(props["accounting-port"]),
		Address:            props["address"],
		AuthenticationPort: parseInt(props["authentication-port"]),
		CalledId:           props["called-id"],
		Certificate:        props["certificate"],
		Comment:            props["comment"],
		Disabled:           parseBool(props[" disabled"]),
		Domain:             props["domain"],
		Protocol:           props["protocol"],
		Realm:              props["realm"],
		Secret:             props["secret"],
		SrcAddress:         props["src-address"],
		Timeout:            parseDuration(props["timeout"]),
	}
	entry.Service = strings.Split(props["service"], ",")
	return entry
}

// GetRadius returns a list of all radius services
func (c *Client) GetRadius() ([]RadiusServer, error) {
	entries := make([]RadiusServer, 0, 8)
	detail, err := c.RunCmd("/radius/print")
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseRadius(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// AddRadius adds a new Radius service
// placeBefore should be empty to just append, else this should be the ID of the entry to which to place this one before
func (c *Client) AddRadius(r RadiusServer, placeBefore string) (string, error) {
	if len(r.Address) == 0 {
		return "", fmt.Errorf("invalid address supplied")
	}
	if len(r.Service) == 0 {
		return "", fmt.Errorf("missing service(s)")
	}
	if r.Protocol != "udp" && r.Protocol != "radsec" {
		return "", fmt.Errorf("invalid protocol")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/radius/add")
	parts = append(parts, fmt.Sprintf("=address=%s", r.Address))
	parts = append(parts, fmt.Sprintf("=protocol=%s", r.Protocol))
	if len(r.CalledId) > 0 {
		parts = append(parts, fmt.Sprintf("=called-id=%s", r.CalledId))
	}
	if r.Protocol == "radsec" && len(r.Certificate) > 0 {
		parts = append(parts, fmt.Sprintf("=certificate=%s", r.Certificate))
	}
	if len(r.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", r.Comment))
	}
	if len(r.Domain) > 0 {
		parts = append(parts, fmt.Sprintf("=domain=%s", r.Domain))
	}
	parts = append(parts, fmt.Sprintf("=disabled=%t", r.Disabled))
	parts = append(parts, fmt.Sprintf("=accounting-backup=%t", r.AccountingBackup))
	if r.AccountingPort > 0 {
		parts = append(parts, fmt.Sprintf("=accounting-port=%d", r.AccountingPort))
	}
	if r.AuthenticationPort > 0 {
		parts = append(parts, fmt.Sprintf("=authentication-port=%d", r.AuthenticationPort))
	}
	if len(r.Realm) > 0 {
		parts = append(parts, fmt.Sprintf("=realm=%s", r.Realm))
	}
	if len(r.Secret) > 0 {
		parts = append(parts, fmt.Sprintf("=secret=%s", r.Secret))
	}
	parts = append(parts, fmt.Sprintf("=service=%s", strings.Join(r.Service, ",")))
	if len(r.SrcAddress) > 0 {
		parts = append(parts, fmt.Sprintf("=src-address=%s", r.SrcAddress))
	}
	if r.Timeout > 0 {
		parts = append(parts, fmt.Sprintf("=timeout=%s", r.Timeout.Round(time.Millisecond).String()))
	}
	if len(placeBefore) > 0 {
		parts = append(parts, fmt.Sprintf("=place-before=%s", placeBefore))
	}
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

// RemoveRadius removes a Radius service by ID. No check is made to see if the item exists.
// Passing an empty ID is a null but successful operation.
func (c *Client) RemoveRadius(id string) error {
	if len(id) > 0 {
		_, err := c.Run("/radius/remove", "=.id="+id)
		if err != nil {
			return err
		}
	}
	return nil
}

// GetAAA returns the User AAA settings
func (c *Client) GetAAA() (AAA, error) {
	var entry AAA
	detail, err := c.RunCmd("/user/aaa/print")
	if err == nil {
		for k, v := range detail.Re[0].Map {
			switch k {
			case "accounting":
				entry.Accounting = parseBool(v)
			case "use-radius":
				entry.UseRadius = parseBool(v)
			case "interim-update":
				entry.InterimUpdate = parseDuration(v)
			case "default-group":
				entry.DefaultGroup = v
			case "exclude-groups":
				entry.ExcludeGroups = strings.Split(v, ",")
			}
		}
	}
	return entry, nil
}

func (c *Client) SetAAA(a AAA) (string, error) {
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/aaa/set")
	parts = append(parts, fmt.Sprintf("=use-radius=%t", a.UseRadius))
	parts = append(parts, fmt.Sprintf("=accounting=%t", a.Accounting))
	if a.InterimUpdate > 0 {
		parts = append(parts, fmt.Sprintf("=interim-update=%s", a.InterimUpdate.Round(time.Millisecond).String()))
	}
	if len(a.DefaultGroup) > 0 {
		parts = append(parts, fmt.Sprintf("=default-group=%s", a.DefaultGroup))
	}
	parts = append(parts, fmt.Sprintf("=exclude-groups=%s", strings.Join(a.ExcludeGroups, ",")))
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}
