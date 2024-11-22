package gotik

import (
	"fmt"
	"strings"
)

type Script struct {
	ID                 string   `json:"id"`
	Comment            string   `json:"comment"`
	DontReqPermissions bool     `json:"dont_req_permissions"`
	Name               string   `json:"name"`
	Owner              string   `json:"owner"`
	Policy             []string `json:"policy"`
	RunCount           int      `json:"run-count"`
	Source             string   `json:"source"`
}

func (e *Script) String() string {
	return fmt.Sprintf("%s %s %s", e.Name, e.Policy, e.Comment)
}

func parseScript(props map[string]string) Script {
	entry := Script{
		ID:                 props[".id"],
		Comment:            props["comment"],
		DontReqPermissions: parseBool(props["dont-require-permissions"]),
		Name:               props["name"],
		Owner:              props["owner"],
		RunCount:           parseInt(props["run-count"]),
		Source:             props["source"],
	}
	entry.Policy = strings.Split(props["policy"], ",")
	return entry
}

// GetScripts returns a list of all scripts
func (c *Client) GetScripts() ([]Script, error) {
	entries := make([]Script, 0, 8)
	detail, err := c.RunCmd("/system/script/print")
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseScript(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// AddScript adds a new Script
func (c *Client) AddScript(s Script) (string, error) {
	if len(s.Name) == 0 {
		return "", fmt.Errorf("invalid name supplied")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/system/script/add")
	parts = append(parts, s.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

// RemoveScript removes a Script service by ID. No check is made to see if the item exists.
// Passing an empty ID is a null but successful operation.
func (c *Client) RemoveScript(id string) error {
	if len(id) > 0 {
		_, err := c.Run("/system/script/remove", "=.id="+id)
		if err != nil {
			return err
		}
	}
	return nil
}

// UpdateScript updates an existing Script.  It will error if the ID or Name are missing or if the
// existing script is not found.
func (c *Client) UpdateScript(s Script) (string, error) {
	if len(s.ID) == 0 {
		return "", fmt.Errorf("missing ID")
	}
	if len(s.Name) == 0 {
		return "", fmt.Errorf("invalid name")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/system/script/set", "=.id="+s.ID)
	parts = append(parts, s.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (s *Script) parter() []string {
	parts := make([]string, 0, 10)
	parts = append(parts, fmt.Sprintf("=name=%s", s.Name))
	if len(s.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", s.Comment))
	}
	parts = append(parts, fmt.Sprintf("=dont-require-permissions=%t", s.DontReqPermissions))
	if len(s.Policy) > 0 {
		parts = append(parts, fmt.Sprintf("=policy=%s", strings.Join(s.Policy, ",")))
	}
	parts = append(parts, fmt.Sprintf("=source=%s", s.Source))
	return parts
}
