package gotik

import (
	"errors"
	"fmt"
	"net"
	"strings"
)

func parsev4Route(props map[string]string) IPv4Route {
	entry := IPv4Route{
		ID:            props[".id"],
		Gateway:       props["gateway"],
		GatewayStatus: props["gateway-status"],
		DstAddress:    props["dst-address"],
		Active:        parseBool(props["active"]),
		Disabled:      parseBool(props["disabled"]),
		Static:        parseBool(props["static"]),
		Connected:     parseBool(props["connected"]),
		Comment:       props["comment"],
		RouteType:     props["type"],
		PrefSrc:       props["pref-src"],
		Mark:          props["routing-mark"],
		Distance:      parseInt(props["distance"]),
		Scope:         parseInt(props["scope"]),
		TargetScope:   parseInt(props["target-scope"]),
	}
	return entry
}

// GetIPv4Routes returns a slice of all routes with optional limiters (ospf, static, connected, disabled)
func (c *Client) GetIPv4Routes(limiters []string) ([]IPv4Route, error) {
	var parms []string

	for _, l := range limiters {
		switch l {
		case "ospf":
			parms = append(parms, "?ospf=true")
		case "static":
			parms = append(parms, "?static=true")
		case "connected":
			parms = append(parms, "?connect=true")
		case "disabled":
			parms = append(parms, "?disabled=true")
		case "enabled":
			parms = append(parms, "?disabled=false")
		case "active":
			parms = append(parms, "?active=true")
		}
	}
	routes := make([]IPv4Route, 0, 1024)
	detail, err := c.RunCmd("/ip/route/print", parms...)
	if err != nil {
		return routes, err
	}
	for _, re := range detail.Re {
		routes = append(routes, parsev4Route(re.Map))
	}
	return routes, nil
}

// FindMatchingIPv4Routes returns a slice of routes where the route gateway falls in the specified subnet
func (c *Client) FindMatchingIPv4Routes(routes []IPv4Route, netToMatch *net.IPNet) (matchingRoutes []IPv4Route, err error) {
	for i := range routes {
		if netToMatch.Contains(net.ParseIP(routes[i].Gateway)) {
			matchingRoutes = append(matchingRoutes, routes[i])
		}
	}

	return
}

// ModifyIPv4Route will modify the specified property of a route
func (c *Client) ModifyIPv4Route(id string, action string) error {
	switch action {
	case "enable":
		_, err := c.Run("/ip/route/set", "=disabled=no", "=.id="+id)
		if err != nil {
			return err
		}
	case "disable":
		_, err := c.Run("/ip/route/set", "=disabled=yes", "=.id="+id)
		if err != nil {
			return err
		}
	case "remove":
		_, err := c.Run("/ip/route/remove", "=.id="+id)
		if err != nil {
			return err
		}
	default:
		return errors.New("route modification action invalid")
	}

	// return nil if all good
	return nil
}

func (r *IPv4Route) String() string {
	a := make([]string, 0, 24)
	f := []rune{' ', ' ', ' ', ' '}
	if r.Disabled {
		f[0] = 'X'
	}
	if r.Active {
		f[1] = 'A'
	}
	if r.Connected {
		f[2] = 'C'
	}
	if r.Static {
		f[2] = 'S'
	}
	a = append(a, string(f), fmt.Sprintf("dst-address=%s", r.DstAddress))
	if len(r.PrefSrc) > 0 {
		a = append(a, fmt.Sprintf("pref-src=%s", r.PrefSrc))
	}
	a = append(a, fmt.Sprintf("gateway=%s", r.Gateway))
	a = append(a, fmt.Sprintf("gateway-status=%s", r.GatewayStatus))
	a = append(a, fmt.Sprintf("distance=%d", r.Distance))
	a = append(a, fmt.Sprintf("scope=%d", r.Scope))
	if len(r.Comment) > 0 {
		a = append(a, fmt.Sprintf("comment=%s", r.Comment))
	}
	return strings.Join(a, " ")
}
