package gotik

import "fmt"

type IPService struct {
	ID          string `json:"id"`
	Disabled    bool   `json:"disabled"`
	Invalid     bool   `json:"invalid"`
	Name        string `json:"name"`
	Port        int    `json:"port"`
	Address     string `json:"address"`
	Certificate string `json:"certificate"`
	TLSVersion  string `json:"tls-version"`
}

func parseIPService(props map[string]string) IPService {
	entry := IPService{
		ID:          props[".id"],
		Disabled:    parseBool(props["disabled"]),
		Invalid:     parseBool(props["invalid"]),
		Name:        props["name"],
		Port:        parseInt(props["port"]),
		Address:     props["address"],
		Certificate: props["certificate"],
		TLSVersion:  props["tls-version"],
	}
	return entry
}

func (c *Client) GetIPServices() ([]IPService, error) {
	entries := make([]IPService, 0)
	detail, err := c.RunCmd("/ip/service/print")
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseIPService(detail.Re[i].Map))
		}
		return entries, nil
	}
	return nil, err
}

// SetIPServiceDisable will enable or disable a service.  Pass either id or name but not both.
func (c *Client) SetIPServiceDisable(id string, disabled bool) error {
	d := "=disabled=true"
	if !disabled {
		d = "=disabled=false"
	}
	_, err := c.Run("/ip/service/set", "=.id="+id, d)
	if err != nil {
		return err
	}
	return nil
}

// SetIPService can set any of the parameters of a service.
// The parameters, id, disabled and address must be supplied.
// Others are optional: port (set to 0 for no change), cert to "" for no change, and tlsVersion to "" for no change.
func (c *Client) SetIPService(id string, disabled bool, port int, address string, cert string, tlsVersion string) error {
	parms := []string{"/ip/service/set", "=.id=" + id, "=disabled=true", "=address=" + address}
	if !disabled {
		parms[2] = "=disabled=false"
	}
	if port > 0 {
		parms = append(parms, fmt.Sprintf("=port=%d", port))
	}
	if len(cert) > 0 {
		parms = append(parms, "=certificate="+cert)
	}
	if len(tlsVersion) > 0 {
		parms = append(parms, "=tls-version="+tlsVersion)
	}
	_, err := c.Run(parms...)
	if err != nil {
		return err
	}
	return nil
}
