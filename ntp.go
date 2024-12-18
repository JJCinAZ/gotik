package gotik

import (
	"fmt"
	"strings"
	"time"
)

// NTP client settings for RouterOS v6 or earlier
type NTPClient6 struct {
	Enabled             bool          `json:"enabled"`
	Servers             []string      `json:"servers"` // may be IPv4 or IPv6
	ServerDNSNames      []string      `json:"server_dns_names"`
	Mode                string        `json:"mode"` // broadcast, manycast, multicast, or unicast
	PollInterval        time.Duration `json:"poll_interval"`
	ActiveServer        string        `json:"active_server"`          // read-only
	LastUpdateFrom      string        `json:"last_update_from"`       // read-only
	TimeSinceLastUpdate time.Duration `json:"time_since_last_update"` // read-only
	LastAdjustment      time.Duration `json:"last_adjustment"`        // read-only
}

// NTP client settings for RouterOS v7 or later
type NTPClient7 struct {
	Enabled        bool          `json:"enabled"`
	Servers        []string      `json:"servers"` // FQDN, ipv4, ipv4@vrf, ipv6, ipv6@vrf, ipv6-linklocal%interface
	Mode           string        `json:"mode"`    // broadcast, manycast, multicast, or unicast
	VRF            string        `json:"vrf"`
	Status         string        `json:"status"`           // read-only
	LastUpdateFrom string        `json:"last_update_from"` // read-only
	SyncedStratum  int           `json:"synced_stratum"`   // read-only
	SystemOffset   time.Duration `json:"system_offset"`    // read-only
}

func parseNTP7(props map[string]string) NTPClient7 {
	entry := NTPClient7{
		Enabled:        parseBool(props["enabled"]),
		Servers:        strings.Split(props["servers"], ","),
		Mode:           props["mode"],
		VRF:            props["vrf"],
		Status:         props["status"],
		LastUpdateFrom: props["synced-server"],
		SyncedStratum:  parseInt(props["synced-stratum"]),
		SystemOffset:   parseDuration(props["system-offset"]),
	}
	return entry
}

func parseNTP6(props map[string]string) NTPClient6 {
	entry := NTPClient6{
		Enabled:        parseBool(props["enabled"]),
		ServerDNSNames: strings.Split(props["server-dns-names"], ","),
		Mode:           props["mode"],
		PollInterval:   parseDuration(props["poll-interval"]),
		LastUpdateFrom: props["last-update-from"],
		LastAdjustment: parseDuration(props["last-adjustment"]),
	}
	if s, found := props["primary-ntp"]; found {
		entry.Servers = append(entry.Servers, s)
	}
	if s, found := props["secondary-ntp"]; found {
		entry.Servers = append(entry.Servers, s)
	}
	return entry
}

// GetNTPClient returns the NTP Client Settings
func (c *Client) GetNTPClient() (any, error) {
	detail, err := c.RunCmd("/system/ntp/client/print")
	if err == nil {
		if c.majorVersion >= 7 {
			return parseNTP7(detail.Re[0].Map), nil
		}
		return parseNTP6(detail.Re[0].Map), nil
	}
	return nil, err
}

// SetNTPClient sets the NTP Client settings.  Must be supplied with an NTPClient6 or NTPClient7 struct.
func (c *Client) SetNTPClient(ntp any) error {
	parts := make([]string, 0, 10)
	parts = append(parts, "/system/ntp/client/set")
	if c.majorVersion >= 7 {
		parts = append(parts, fmt.Sprintf("=enabled=%t", ntp.(NTPClient7).Enabled))
		if len(ntp.(NTPClient7).Servers) > 0 {
			parts = append(parts, fmt.Sprintf("=servers=%s", strings.Join(ntp.(NTPClient7).Servers, ",")))
		}
		if ntp.(NTPClient7).Mode != "" {
			parts = append(parts, fmt.Sprintf("=mode=%s", ntp.(NTPClient7).Mode))
		}
		if ntp.(NTPClient7).VRF != "" {
			parts = append(parts, fmt.Sprintf("=vrf=%s", ntp.(NTPClient7).VRF))
		}
	} else {
		parts = append(parts, fmt.Sprintf("=enabled=%t", ntp.(NTPClient6).Enabled))
		if len(ntp.(NTPClient6).Servers) > 0 {
			parts = append(parts, fmt.Sprintf("=primary-ntp=%s", ntp.(NTPClient6).Servers[0]))
			if len(ntp.(NTPClient6).Servers) > 1 {
				parts = append(parts, fmt.Sprintf("=secondary-ntp=%s", ntp.(NTPClient6).Servers[1]))
			} else {
				parts = append(parts, "=secondary-ntp=")
			}
		}
		if ntp.(NTPClient6).ServerDNSNames != nil && len(ntp.(NTPClient6).ServerDNSNames) > 0 {
			parts = append(parts, fmt.Sprintf("=server-dns-names=%s", strings.Join(ntp.(NTPClient6).ServerDNSNames, ",")))
		}
	}
	_, err := c.Run(parts...)
	return err
}
