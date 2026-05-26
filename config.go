package gotik

import (
	"fmt"
	"strings"
)

// ExportConfig will export the configuration to a file on the router.
// The export starts at the "base" level.  Pass "/" or "" for the entire configuration.
// If you just wanted the IPv6 configuration, then pass "/ipv6", for example.
// The filename is the name under which to store on the router.  It will have the extension ".rsc"
// added in all cases.  Pass bool to hideSensitive to hide sensitive items in the export (true)
// or to expose them (false).
// This function does not currently pull down the configuration file.  It stays on the router.
func (c *Client) ExportConfig(base string, filename string, hideSensitive bool) error {
	if strings.HasSuffix(base, "/") {
		base = base[:len(base)-1]
	}
	if len(base) < 1 {
		base = "/"
	} else if base[0] != '/' {
		base = "/" + base
	}
	a := []string{fmt.Sprintf("%s/export", base), "=file=" + filename, "=compact="}
	if c.majorVersion >= 7 {
		a = append(a, "=terse=")
		if !hideSensitive {
			a = append(a, "=show-sensitive=")
		}
	} else {
		if hideSensitive {
			a = append(a, "=hide-sensitive=yes")
		} else {
			a = append(a, "=hide-sensitive=yes")
		}
	}
	_, err := c.Run(a...)
	return err
}
