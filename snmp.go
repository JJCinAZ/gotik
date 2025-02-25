package gotik

import (
	"fmt"
	"strings"
)

func (c *SNMPCommunity) String() string {
	flags := ' '
	if c.Disabled {
		flags = 'X'
	}
	return fmt.Sprintf("%c %s addresses=%s community=%s read=%t write=%t %s", flags,
		c.Name, c.Addresses, c.Name, c.ReadAccess, c.WriteAccess, c.Comment)
}

func parseSNMPCommunity(props map[string]string) SNMPCommunity {
	entry := SNMPCommunity{
		ID:                     props[".id"],
		Disabled:               parseBool(props["disabled"]),
		Default:                parseBool(props["default"]),
		Comment:                props["comment"],
		Name:                   props["name"],
		ReadAccess:             parseBool(props["read-access"]),
		WriteAccess:            parseBool(props["write-access"]),
		AuthenticationProtocol: props["authentication-protocol"],
		AuthenticationPassword: props["authentication-password"],
		EncryptionProtocol:     props["encryption-protocol"],
		EncryptionPassword:     props["encryption-password"],
		Security:               props["security"],
	}
	entry.Addresses = strings.Split(props["addresses"], ",")
	return entry
}

// GetSNMPCommunities returns a list of all SNMP communities
func (c *Client) GetSNMPCommunities() ([]SNMPCommunity, error) {
	entries := make([]SNMPCommunity, 0, 8)
	detail, err := c.RunCmd("/snmp/community/print")
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseSNMPCommunity(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// AddSNMPCommunity adds a new SNMP Community
// The Default field is a read-only field and ignored for adds.  The first community is always the default
// and always exists on a router (it cannot be removed).
func (c *Client) AddSNMPCommunity(community SNMPCommunity) (string, error) {
	parts := make([]string, 0)
	parts = append(parts, "/snmp/community/add")
	parts = append(parts, community.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

// UpdateSNMPCommunity will update an existing SNMP Community
// The ID field must be set to the ID of the community to update.
// The Default field is a read-only field and ignored for updates.
func (c *Client) UpdateSNMPCommunity(community SNMPCommunity) (string, error) {
	if len(community.ID) == 0 {
		return "", fmt.Errorf("missing ID")
	}
	if len(community.Name) == 0 {
		return "", fmt.Errorf("invalid community Name")
	}
	parts := make([]string, 0)
	parts = append(parts, "/snmp/community/set", "=.id="+community.ID)
	parts = append(parts, community.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (community *SNMPCommunity) parter() []string {
	var parts []string
	parts = append(parts, fmt.Sprintf("=disabled=%t", community.Disabled))
	if len(community.Addresses) > 0 {
		parts = append(parts, fmt.Sprintf("=addresses=%s", strings.Join(community.Addresses, ",")))
	}
	if len(community.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", community.Comment))
	}
	if len(community.Name) > 0 {
		parts = append(parts, fmt.Sprintf("=name=%s", community.Name))
	}
	if len(community.AuthenticationProtocol) > 0 {
		parts = append(parts, fmt.Sprintf("=authentication-protocol=%s", community.AuthenticationProtocol))
	}
	if len(community.AuthenticationPassword) > 0 {
		parts = append(parts, fmt.Sprintf("=authentication-password=%s", community.AuthenticationPassword))
	}
	if len(community.EncryptionProtocol) > 0 {
		parts = append(parts, fmt.Sprintf("=encryption-protocol=%s", community.EncryptionProtocol))
	}
	if len(community.EncryptionPassword) > 0 {
		parts = append(parts, fmt.Sprintf("=encryption-password=%s", community.EncryptionPassword))
	}
	if len(community.Security) > 0 {
		parts = append(parts, fmt.Sprintf("=security=%s", community.Security))
	} else {
		parts = append(parts, "=security=none")
	}
	parts = append(parts, fmt.Sprintf("=read-access=%t", community.ReadAccess))
	parts = append(parts, fmt.Sprintf("=write-access=%t", community.WriteAccess))
	return parts
}

// RemoveSNMPCommunity removes a SNMP Community by ID. No check is made to see if the item exists.
// Passing an empty ID is a null but successful operation.
// Note that the default community cannot be removed.
func (c *Client) RemoveSNMPCommunity(id string) error {
	if len(id) > 0 {
		_, err := c.Run("/snmp/community/remove", "=.id="+id)
		if err != nil {
			return err
		}
	}
	return nil
}
