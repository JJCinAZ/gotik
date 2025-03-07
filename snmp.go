package gotik

import (
	"fmt"
	"strings"
)

type SNMPCommunity struct {
	ID                     string   `tik:".id"`
	Disabled               bool     `tik:"disabled"`
	Default                bool     `tik:"default"`
	Comment                string   `tik:"comment"`
	Name                   string   `tik:"name"`
	Addresses              []string `tik:"addresses"` // list of CIDRs (v4 or v6)
	ReadAccess             bool     `tik:"read-access"`
	WriteAccess            bool     `tik:"write-access"`
	AuthenticationProtocol string   `tik:"authentication-protocol"` // MD5 or SHA1 (or blank)
	AuthenticationPassword string   `tik:"authentication-password"`
	EncryptionProtocol     string   `tik:"encryption-protocol"` // AES or DES (or blank)
	EncryptionPassword     string   `tik:"encryption-password"`
	Security               string   `tik:"security"` // none, authorized, or private
}

type SNMP struct {
	Contact        string   `tik:"contact"`
	Enabled        bool     `tik:"enabled"`
	EngineId       string   `tik:"engine-id"`        // For 7.10 or newer, this is the hexidecimal string and is read-only
	EngineIdSuffix string   `tik:"engine-id-suffix"` // only for 7.10 or newer
	Location       string   `tik:"location"`
	SrcAddress     string   `tik:"src-address"` // defaults to "::"
	TrapCommunity  string   `tik:"trap-community"`
	TrapGenerators []string `tik:"trap-generators"` // possible values: 'interfaces', 'start-trap', 'temp-exception'
	TrapInterfaces []string `tik:"trap-interfaces"` // list of interface names separated by commas or "all"
	TrapTarget     string   `tik:"trap-target"`     // list of IP addresses separated by commas (can be IPv4 or IPv6)
	TrapVersion    string   `tik:"trap-version"`    // 1, 2, or 3
	VRF            string   `tik:"vrf"`             // Only 7.3 or newer
}

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
	parts = append(parts, community.parter(c)...)
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
	parts = append(parts, community.parter(c)...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (community *SNMPCommunity) parter(c *Client) []string {
	var parts []string
	if c.majorVersion > 6 || (c.majorVersion == 6 && c.minorVersion >= 46) {
		parts = append(parts, fmt.Sprintf("=disabled=%t", community.Disabled))
		parts = append(parts, fmt.Sprintf("=comment=%s", community.Comment))
	}
	if len(community.Addresses) > 0 {
		parts = append(parts, fmt.Sprintf("=addresses=%s", strings.Join(community.Addresses, ",")))
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

func parseSNMP(props map[string]string) SNMP {
	entry := SNMP{
		Contact:        props["contact"],
		Enabled:        parseBool(props["enabled"]),
		EngineId:       props["engine-id"],
		EngineIdSuffix: props["engine-id-suffix"],
		Location:       props["location"],
		SrcAddress:     props["src-address"],
		TrapCommunity:  props["trap-community"],
		TrapGenerators: strings.Split(props["trap-generators"], ","),
		TrapInterfaces: strings.Split(props["trap-interfaces"], ","),
		TrapTarget:     props["trap-target"],
		TrapVersion:    props["trap-version"],
		VRF:            props["vrf"],
	}
	return entry
}

// GetSNMP returns the DNS settings
func (c *Client) GetSNMP() (SNMP, error) {
	detail, err := c.RunCmd("/snmp/print")
	if err == nil {
		return parseSNMP(detail.Re[0].Map), nil
	}
	return SNMP{}, err
}

// SetSNMP sets the SNMP settings.  Note that any communities must have been added already.
func (c *Client) SetSNMP(s SNMP) error {
	parts := make([]string, 0, 10)
	parts = append(parts, "/snmp/set")
	if s.Enabled {
		parts = append(parts, "=enabled=yes")
	} else {
		parts = append(parts, "=enabled=no")
	}
	parts = append(parts, fmt.Sprintf("=contact=%s", s.Contact))
	parts = append(parts, fmt.Sprintf("=location=%s", s.Location))
	parts = append(parts, fmt.Sprintf("=trap-version=%s", s.TrapVersion))
	parts = append(parts, fmt.Sprintf("=trap-community=%s", s.TrapCommunity))
	parts = append(parts, fmt.Sprintf("=trap-target=%s", s.TrapTarget))
	parts = append(parts, fmt.Sprintf("=trap-generators=%s", strings.Join(s.TrapGenerators, ",")))
	if c.majorVersion > 6 || (c.majorVersion == 6 && c.minorVersion >= 44) {
		if s.SrcAddress == "" || s.SrcAddress == "::" {
			// Once src-address is set, it cannot be set back to blank, so we either have to pick
			// 0.0.0.0 or :: based on whether ipv6 is enabled
			hasIPv6, err := c.IsPackageEnabled("ipv6")
			if err != nil {
				return err
			}
			if hasIPv6 {
				s.SrcAddress = "::"
			} else {
				s.SrcAddress = "0.0.0.0"
			}
		}
		parts = append(parts, fmt.Sprintf("=src-address=%s", s.SrcAddress))
	}
	if c.majorVersion > 7 || (c.majorVersion == 7 && c.minorVersion >= 3) {
		if len(s.VRF) > 0 {
			parts = append(parts, fmt.Sprintf("=vrf=%s", s.VRF))
		} else {
			parts = append(parts, "=vrf=main")
		}
		if len(s.TrapInterfaces) == 0 {
			parts = append(parts, "=trap-interfaces=all")
		} else {
			parts = append(parts, fmt.Sprintf("=trap-interfaces=%s", strings.Join(s.TrapInterfaces, ",")))
		}
	}
	if c.majorVersion > 7 || (c.majorVersion == 7 && c.minorVersion >= 10) {
		parts = append(parts, fmt.Sprintf("=engine-id-suffix=%s", s.EngineIdSuffix))
	} else {
		parts = append(parts, fmt.Sprintf("=engine-id=%s", s.EngineId))
	}
	_, err := c.Run(parts...)
	return err
}
