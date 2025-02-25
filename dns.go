package gotik

import (
	"fmt"
	"strings"
	"time"
)

type DNS struct {
	AllowRemoteRequests      bool          `json:"allow_remote_requests"`
	CacheSize                int           `json:"cache_size"` // KiB
	CacheUsed                int           `json:"cache_used"` // KiB
	CacheMaxTTL              time.Duration `json:"cache_max_ttl"`
	MaxConcurrentQueries     int           `json:"max_concurrent_queries"`
	MaxConcurrentTCPSessions int           `json:"max_concurrent_tcp_sessions"`
	MaxUDPPacketSize         int           `json:"max_udp_packet_size"`
	QueryServerTimeout       time.Duration `json:"query_server_timeout"`
	QueryTotalTimeout        time.Duration `json:"query_total_timeout"`
	UseDOHServer             string        `json:"use_doh_server"`
	VerifyDOHCert            bool          `json:"verify_doh_cert"`
	Servers                  []string      `json:"servers"`
	DynamicServers           []string      `json:"dynamic_servers"`
}

func parseDNS(props map[string]string) DNS {
	entry := DNS{
		Servers:                  strings.Split(props["servers"], ","),
		DynamicServers:           strings.Split(props["dynamic-servers"], ","),
		UseDOHServer:             props["use-doh-server"],
		VerifyDOHCert:            parseBool(props["verify-doh-cert"]),
		AllowRemoteRequests:      parseBool(props["allow-remote-requests"]),
		MaxUDPPacketSize:         parseInt(props["max-udp-packet-size"]),
		QueryServerTimeout:       parseDuration(props["query-server-timeout"]),
		QueryTotalTimeout:        parseDuration(props["query-total-timeout"]),
		MaxConcurrentQueries:     parseInt(props["max-concurrent-queries"]),
		MaxConcurrentTCPSessions: parseInt(props["max-concurrent-tcp-sessions"]),
		CacheSize:                parseInt(props["cache-size"]),
		CacheUsed:                parseInt(props["cache-used"]),
		CacheMaxTTL:              parseDuration(props["cache-max-ttl"]),
	}
	return entry
}

// GetDNS returns the DNS settings
func (c *Client) GetDNS() (DNS, error) {
	detail, err := c.RunCmd("/ip/dns/print")
	if err == nil {
		return parseDNS(detail.Re[0].Map), nil
	}
	return DNS{}, err
}

// FlushDNS flushes the DNS cache
func (c *Client) FlushDNS() error {
	_, err := c.Run("/ip/dns/cache/flush")
	return err
}

// SetDNS sets the DNS settings
func (c *Client) SetDNS(d DNS) error {
	parts := make([]string, 0, 10)
	parts = append(parts, "/ip/dns/set")
	if d.CacheSize > 0 {
		parts = append(parts, fmt.Sprintf("=cache-size=%d", d.CacheSize))
	}
	if d.CacheMaxTTL > 0 {
		parts = append(parts, fmt.Sprintf("=cache-max-ttl=%s", d.CacheMaxTTL.Round(time.Millisecond).String()))
	}
	if d.MaxConcurrentQueries > 0 {
		parts = append(parts, fmt.Sprintf("=max-concurrent-queries=%d", d.MaxConcurrentQueries))
	}
	if d.MaxConcurrentTCPSessions > 0 {
		parts = append(parts, fmt.Sprintf("=max-concurrent-tcp-sessions=%d", d.MaxConcurrentTCPSessions))
	}
	if d.MaxUDPPacketSize > 0 {
		parts = append(parts, fmt.Sprintf("=max-udp-packet-size=%d", d.MaxUDPPacketSize))
	}
	if d.QueryServerTimeout > 0 {
		parts = append(parts, fmt.Sprintf("=query-server-timeout=%s", d.QueryServerTimeout.Round(time.Millisecond).String()))
	}
	if d.QueryTotalTimeout > 0 {
		parts = append(parts, fmt.Sprintf("=query-total-timeout=%s", d.QueryTotalTimeout.Round(time.Millisecond).String()))
	}
	if c.majorVersion > 6 || (c.majorVersion == 6 && c.minorVersion >= 49) {
		parts = append(parts, fmt.Sprintf("=use-doh-server=%s", d.UseDOHServer))
		parts = append(parts, fmt.Sprintf("=verify-doh-cert=%t", d.VerifyDOHCert))
	}
	parts = append(parts, fmt.Sprintf("=allow-remote-requests=%t", d.AllowRemoteRequests))
	parts = append(parts, fmt.Sprintf("=servers=%s", strings.Join(d.Servers, ",")))
	_, err := c.Run(parts...)
	return err
}
