package gotik

import (
	"fmt"
)

func parseUser(props map[string]string) User {
	entry := User{
		ID:        props[".id"],
		Name:      props["name"],
		Group:     props["group"],
		Comment:   props["comment"],
		Disabled:  parseBool(props["disabled"]),
		LastLogin: props["last-logged-in"],
		Address:   props["address"],
	}
	return entry
}

func (c *Client) userPrint(parms ...string) ([]User, error) {
	entries := make([]User, 0)
	detail, err := c.RunCmd("/user/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parseUser(detail.Re[i].Map))
		}
	} else {
		entries = nil
	}
	return entries, err
}

// Returns a User by name
func (c *Client) GetUserByName(name string) (User, error) {
	a, err := c.userPrint("?name=" + name)
	if err == nil {
		if len(a) > 0 {
			return a[0], nil
		}
		err = ErrNotFound
	}
	return User{}, err
}

// Returns all Users
func (c *Client) GetUsers() ([]User, error) {
	return c.userPrint()
}

// Add or update a User.
// If a user with the same name already exists, it will be updated.
// Password will not be set for user
func (c *Client) AddUser(user User) (string, error) {
	if len(user.Name) == 0 {
		return "", fmt.Errorf("invalid User")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/add")
	parts = append(parts, user.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	if apiErr, okay := err.(*DeviceError); okay {
		if msg, found := apiErr.Sentence.Map["message"]; found {
			if msg == "failure: user with the same name already exists" {
				goto UPDATE
			}
		}
	}
	return "", err
UPDATE:
	var existing User
	if existing, err = c.GetUserByName(user.Name); err == nil {
		user.ID = existing.ID
		return c.UpdateUser(user)
	}
	return "", err
}

/*
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
*/

// UpdateUser updates an existing user.  ID must be valid.
func (c *Client) UpdateUser(user User) (string, error) {
	if len(user.ID) == 0 {
		return "", fmt.Errorf("missing ID")
	}
	if len(user.Name) == 0 {
		return "", fmt.Errorf("invalid User")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/set", "=.id="+user.ID)
	parts = append(parts, user.parter()...)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (c *Client) UpdateUserPasswordByID(ID, password string) (string, error) {
	if len(ID) == 0 {
		return "", fmt.Errorf("missing ID")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/set", "=.id="+ID)
	parts = append(parts, "=password="+password)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (c *Client) UpdateUserPasswordByName(username, password string) (string, error) {
	if len(username) == 0 {
		return "", fmt.Errorf("invalid User")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/user/set", "=name="+username)
	parts = append(parts, "=password="+password)
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}

func (u *User) parter() []string {
	parts := make([]string, 0, 10)
	parts = append(parts, fmt.Sprintf("=name=%s", u.Name))
	parts = append(parts, fmt.Sprintf("=disabled=%t", u.Disabled))
	if len(u.Address) > 0 {
		parts = append(parts, fmt.Sprintf("=address=%s", u.Address))
	}
	if len(u.Group) > 0 {
		parts = append(parts, fmt.Sprintf("=group=%s", u.Group))
	}
	if len(u.Comment) > 0 {
		parts = append(parts, fmt.Sprintf("=comment=%s", u.Comment))
	}
	return parts
}
