package gotik

import "fmt"

func parsePPPoEServer(props map[string]string) PPPoEServer {
	entry := PPPoEServer{
		ID:             props[".id"],
		Disabled:       parseBool(props["disabled"]),
		Interface:      props["interface"],
		ServiceName:    props["service-name"],
		MaxMTU:         parseInt(props["max-mtu"]),
		MaxMRU:         parseInt(props["max-mru"]),
		MRRU:           parseInt(props["mrru"]),
		Authentication: props["authentication"],
		KeepAlive:      parseInt(props["keepalive-timeout"]),
		SingleSess:     parseBool(props["one-session-per-host"]),
		MaxSessions:    props["max-sessions"],
		DefaultProfile: props["default-profile"],
		PadoDelay:      parseInt(props["pado-delay"]),
	}
	return entry
}

func (c *Client) pppoeServerPrint(parms ...string) ([]PPPoEServer, error) {
	entries := make([]PPPoEServer, 0)
	detail, err := c.RunCmd("/interface/pppoe-server/server/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parsePPPoEServer(detail.Re[i].Map))
		}
	} else {
		entries = nil
	}
	return entries, err
}

// GetPPPoEServers returns a list of all PPPoE Servers on a particular interface or all servers if intf is blank
func (c *Client) GetPPPoEServers(intf string) ([]PPPoEServer, error) {
	if len(intf) > 0 {
		return c.pppoeServerPrint("?=interface=" + intf)
	} else {
		return c.pppoeServerPrint()
	}
}

// RemovePPPoEServer removes a PPPoE server by ID
func (c *Client) RemovePPPoEServer(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("missing ID")
	}
	_, err := c.Run("/interface/pppoe-server/server/remove", "=.id="+id)
	return err
}
