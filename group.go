package gotik

import (
	"errors"
	"fmt"
	"strings"
)

var groupPolicies = []string{
	"local", "telnet", "ssh", "ftp", "reboot", "read", "write", "policy", "test", "winbox", "password", "web", "sniff",
	"sensitive", "api", "romon", "rest-api",
}

type Group struct {
	ID     string
	Name   string
	Skin   string
	Policy map[string]bool
}

func parseGroup(props map[string]string) Group {
	g := Group{
		ID:     props[".id"],
		Name:   props["name"],
		Skin:   props["skin"],
		Policy: make(map[string]bool),
	}
	a := strings.Split(props["policy"], ",")
	for _, policy := range a {
		if len(policy) == 0 {
			continue
		}
		if policy[0] == '!' {
			g.Policy[policy[1:]] = false
			continue
		}
		g.Policy[policy] = true
	}
	return g
}

func (c *Client) groupPrint(parms ...string) ([]Group, error) {
	entries := make([]Group, 0)
	detail, err := c.RunCmd("/user/group/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseGroup(detail.Re[i].Map))
		}
	} else {
		entries = nil
	}
	return entries, err
}

func (c *Client) GetGroups() ([]Group, error) {
	return c.groupPrint()
}

func (c *Client) GetGroupByName(name string) (Group, error) {
	a, err := c.groupPrint("?name=" + name)
	if err == nil {
		if len(a) > 0 {
			return a[0], nil
		}
		err = ErrNotFound
	}
	return Group{}, err
}

func (c *Client) AddGroup(g Group) (string, error) {
	if len(g.Name) == 0 {
		return "", fmt.Errorf("invalid Group")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/group/add")
	parts = append(parts, fmt.Sprintf("=name=%s", g.Name))
	if len(g.Skin) > 0 {
		parts = append(parts, fmt.Sprintf("=skin=%s", g.Skin))
	}
	parts = append(parts, fmt.Sprintf("=policy=%s", g.assemblePolicies()))
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	var apiErr *DeviceError
	if errors.As(err, &apiErr) {
		if msg, found := apiErr.Sentence.Map["message"]; found {
			if msg == "failure: group with the same name already exists" {
				goto UPDATE
			}
		}
	}
	return "", err
UPDATE:
	var existing Group
	if existing, err = c.GetGroupByName(g.Name); err == nil {
		g.ID = existing.ID
		return c.UpdateGroup(g)
	}
	return "", err
}

func (c *Client) UpdateGroup(g Group) (string, error) {
	if len(g.ID) == 0 {
		return "", fmt.Errorf("missing ID")
	}
	if len(g.Name) == 0 {
		return "", fmt.Errorf("invalid Group")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/group/set", "=.id="+g.ID)
	parts = append(parts, fmt.Sprintf("=name=%s", g.Name))
	parts = append(parts, fmt.Sprintf("=skin=%s", g.Skin))
	parts = append(parts, fmt.Sprintf("=policy=%s", g.assemblePolicies()))
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (g Group) assemblePolicies() string {
	list := make([]string, 0, 10)
	for k, v := range g.Policy {
		if v {
			list = append(list, k)
		} else {
			list = append(list, "!"+k)
		}
	}
	return strings.Join(list, ",")
}

func (c *Client) RemoveGroupByName(name string) error {
	var (
		existing Group
		err      error
	)
	if existing, err = c.GetGroupByName(name); err == nil {
		return c.RemoveGroup(existing.ID)
	}
	return err
}

func (c *Client) RemoveGroup(id string) error {
	if len(id) == 0 {
		return fmt.Errorf("missing ID")
	}
	_, err := c.Run("/user/group/remove", "=.id="+id)
	return err
}
