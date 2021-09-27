package gotik

import (
	"fmt"
)

func parsePPPSecret(props map[string]string) PPPSecret {
	entry := PPPSecret{
		ID:            props[".id"],
		Name:          props["name"],
		CallerID:      props["caller-id"],
		Comment:       props["comment"],
		Disabled:      parseBool(props["disabled"]),
		LimitBytesIn:  parseInt(props["limit-bytes-in"]),
		LimitBytesOut: parseInt(props["limit-bytes-out"]),
		LocalAddress:  props["local-address"],
		Password:      props["password"],
		Profile:       props["profile"],
		RemoteAddress: props["remote-address"],
		Routes:        props["routes"],
		Service:       props["service"],
	}
	return entry
}

func (c *Client) pppSecretPrint(parms ...string) ([]PPPSecret, error) {
	entries := make([]PPPSecret, 0)
	detail, err := c.RunCmd("/ppp/secret/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parsePPPSecret(detail.Re[i].Map))
		}
	} else {
		entries = nil
	}
	return entries, err
}

// Returns a PPP Secret by name
func (c *Client) GetPPPSecretByName(name string) (PPPSecret, error) {
	a, err := c.pppSecretPrint("?name=" + name)
	if err == nil {
		if len(a) > 0 {
			return a[0], nil
		}
		err = ErrNotFound
	}
	return PPPSecret{}, err
}

// Returns all PPP secrets
func (c *Client) GetPPPSecrets() ([]PPPSecret, error) {
	return c.pppSecretPrint()
}

// Add or update a PPP Secret.
// If a secret with the same name already exists, it will be updated.
func (c *Client) AddPPPSecret(secret PPPSecret) (string, error) {
	if len(secret.Name) == 0 {
		return "", fmt.Errorf("invalid PPP Secret")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/ppp/secret/add")
	parts = append(parts, secret.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	if apiErr, okay := err.(*DeviceError); okay {
		if msg, found := apiErr.Sentence.Map["message"]; found {
			if msg == "failure: secret with the same name already exists" {
				goto UPDATE
			}
		}
	}
	return "", err
UPDATE:
	var existing PPPSecret
	if existing, err = c.GetPPPSecretByName(secret.Name); err == nil {
		secret.ID = existing.ID
		return c.UpdatePPPSecret(secret)
	}
	return "", err
}

// Remove a PPP Secret by Name
func (c *Client) RemovePPPSecretByName(name string) error {
	var (
		existing PPPSecret
		err      error
	)
	if existing, err = c.GetPPPSecretByName(name); err == nil {
		return c.RemovePPPSecret(existing.ID)
	}
	return err
}

// Remove a PPP Secret.
func (c *Client) RemovePPPSecret(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("missing ID")
	}
	_, err := c.Run("/ppp/secret/remove", "=.id="+id)
	return err
}

// Update a PPP Secret.
func (c *Client) UpdatePPPSecret(secret PPPSecret) (string, error) {
	if len(secret.ID) == 0 {
		return "", fmt.Errorf("missing ID")
	}
	if len(secret.Name) == 0 {
		return "", fmt.Errorf("invalid PPP Secret")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/ppp/secret/set", "=.id="+secret.ID)
	parts = append(parts, secret.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (secret *PPPSecret) parter() []string {
	parts := make([]string, 0, 10)
	parts = append(parts, fmt.Sprintf("=name=%s", secret.Name))
	parts = append(parts, fmt.Sprintf("=disabled=%t", secret.Disabled))
	if len(secret.Password) > 0 {
		parts = append(parts, fmt.Sprintf("=password=%s", secret.Password))
	}
	if len(secret.Profile) > 0 {
		parts = append(parts, fmt.Sprintf("=profile=%s", secret.Profile))
	}
	if len(secret.Service) > 0 {
		parts = append(parts, fmt.Sprintf("=service=%s", secret.Service))
	}
	if len(secret.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", secret.Comment))
	}
	if len(secret.Routes) > 0 {
		parts = append(parts, fmt.Sprintf("=routes=%s", secret.Routes))
	}
	if len(secret.CallerID) > 0 {
		parts = append(parts, fmt.Sprintf("=caller-id=%s", secret.CallerID))
	}
	if len(secret.LocalAddress) > 0 {
		parts = append(parts, fmt.Sprintf("=local-address=%s", secret.LocalAddress))
	}
	if len(secret.RemoteAddress) > 0 {
		parts = append(parts, fmt.Sprintf("=remote-address=%s", secret.RemoteAddress))
	}
	if secret.LimitBytesIn > 0 {
		parts = append(parts, fmt.Sprintf("=limit-bytes-in=%d", secret.LimitBytesIn))
	}
	if secret.LimitBytesOut > 0 {
		parts = append(parts, fmt.Sprintf("=limit-bytes-out=%d", secret.LimitBytesOut))
	}
	return parts
}

func parsePPPActive(props map[string]string) PPPActive {
	entry := PPPActive{
		ID:            props[".id"],
		Name:          props["name"],
		Address:       props["address"],
		CallerID:      props["caller-id"],
		Radius:        parseBool(props["radius"]),
		LimitBytesIn:  parseInt(props["limit-bytes-in"]),
		LimitBytesOut: parseInt(props["limit-bytes-out"]),
		Service:       props["service"],
		Encoding:      props["encoding"],
		Uptime:        parseDuration(props["uptime"]),
		SessionID:     parseHex(props["session-id"]),
	}
	return entry
}

func (c *Client) pppActivePrint(parms ...string) ([]PPPActive, error) {
	entries := make([]PPPActive, 0, 32)
	detail, err := c.RunCmd("/ppp/active/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parsePPPActive(detail.Re[i].Map))
		}
	} else {
		entries = nil
	}
	return entries, err
}

// Returns all PPP Active connections
func (c *Client) GetPPPActiveConnections() ([]PPPActive, error) {
	return c.pppActivePrint()
}

// Returns a specific PPP Active connection by name
func (c *Client) GetPPPActiveConnectionByName(name string) (PPPActive, error) {
	a, err := c.pppActivePrint("?name=" + name)
	if err == nil {
		if len(a) > 0 {
			return a[0], nil
		}
		err = ErrNotFound
	}
	return PPPActive{}, err
}
