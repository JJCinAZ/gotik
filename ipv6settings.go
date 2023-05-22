package gotik

import (
	"errors"
)

type IPv6Settings struct {
	Forward                    bool
	AcceptRedirects            string
	AcceptRouterAdvertisements string
	MaxNeighborEntries         int
}

func parsev6Settings(props map[string]string) IPv6Settings {
	entry := IPv6Settings{
		Forward:                    parseBool(props["forward"]),
		AcceptRedirects:            props["accept-redirects"],
		AcceptRouterAdvertisements: props["accept-router-advertisements"],
		MaxNeighborEntries:         parseInt(props["max-neighbor-entries"]),
	}
	return entry
}

// GetIPv6Settings returns the current IPv6 settings
func (c *Client) GetIPv6Settings() (settings IPv6Settings, err error) {
	var detail *Reply
	if detail, err = c.Run("/ipv6/settings/print"); err != nil {
		return
	}

	// we only want one result
	switch len(detail.Re) {
	case 0:
		err = ErrNotFound
	case 1:
		settings = parsev6Settings(detail.Re[0].Map)
	default:
		err = errors.New("unexpected return")
	}
	return
}
