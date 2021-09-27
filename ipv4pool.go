package gotik

import "fmt"

func parsev4Pool(props map[string]string) IPv4Pool {
	entry := IPv4Pool{
		ID:       props[".id"],
		Name:     props["name"],
		Ranges:   props["ranges"],
		NextPool: props["next-pool"],
	}
	return entry
}

func (c *Client) ipPoolPrint(parms ...string) ([]IPv4Pool, error) {
	entries := make([]IPv4Pool, 0, 8)
	detail, err := c.RunCmd("/ip/pool/print", parms...)
	if err == nil {
		for i := range detail.Re {
			entries = append(entries, parsev4Pool(detail.Re[i].Map))
		}
	}
	return entries, nil
}

// Returns a single Pool by name
func (c *Client) GetIPv4Pool(name string) ([]IPv4Pool, error) {
	return c.ipPoolPrint("?=name=" + name)
}

// Returns a list of all Pools
func (c *Client) GetIPv4Pools() ([]IPv4Pool, error) {
	return c.ipPoolPrint()
}

// Add a new IPv4 address Pool
func (c *Client) AddIPv4Pool(pool IPv4Pool) (string, error) {
	if len(pool.Name) == 0 || len(pool.Ranges) == 0 {
		return "", fmt.Errorf("invalid pool supplied")
	}
	parts := make([]string, 0, 10)
	parts = append(parts, "/ip/pool/add")
	parts = append(parts, fmt.Sprintf("=name=%s", pool.Name))
	parts = append(parts, fmt.Sprintf("=ranges=%s", pool.Ranges))
	if len(pool.NextPool) > 0 {
		parts = append(parts, fmt.Sprintf("=next-pool=%s", pool.NextPool))
	}
	reply, err := c.Run(parts...)
	if err == nil {
		return reply.Done.Map["ret"], nil
	}
	return "", err
}
