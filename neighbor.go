package gotik

import "errors"

// GetNeighborInterface gets the neighbor discovery interface id by interface
func (c *Client) GetNeighborInterface(iface string) (NeighborInterface, error) {
	// get id of neighbor interface
	detail, err := c.Run("/ip/neighbor/discovery/print", "?name="+iface)
	if err != nil {
		return NeighborInterface{}, err
	}

	// only return the id if there is one result
	switch len(detail.Re) {
	case 0:
		return NeighborInterface{}, errors.New("no neighbor interface - expecting one")
	case 1:
		return NeighborInterface{ID: detail.Re[0].Map[".id"], Discover: detail.Re[0].Map["discover"]}, nil
	default:
		return NeighborInterface{}, errors.New("more than one neighbor interface - only expecting one")
	}
}

// ModifyNeighbor can be used to enable or disable the neighbor discovery interface
func (c *Client) ModifyNeighbor(id string, action string) error {
	var disabled string

	// change disabled sentence based on action
	switch action {
	case "enable":
		disabled = "=discover=true"
	case "disable":
		disabled = "=discover=false"
	default:
		return errors.New("modify neighbor action invalid")
	}

	// set discovery property
	_, err := c.Run("/ip/neighbor/discovery/set", "=.id="+id, disabled)
	if err != nil {
		return err
	}

	// return nil if all good
	return nil
}
