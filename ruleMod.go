package gotik

import (
	"errors"
	"fmt"
	"reflect"
	"strconv"
	"time"
)

// checks for a rule on the router conn
// returns a slice of matching rule ids and the rule location
func (c *Client) GetIDS(in interface{}) ([]string, string, error) {
	var ruleLocation string
	var ids []string

	// get interface reflection information
	s := reflect.ValueOf(in).Elem()
	t := s.Type()

	// iterate through all fields to find router location
	for i := 0; i < s.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tik")

		if field.Name == "RouterLocation" {
			ruleLocation = tag
			break
		}
	}

	// check for rule if location is present
	if len(ruleLocation) > 0 {
		// check for matching rule
		// This builds a command like:
		// /ip/firewall/filter/print ?action=reject ?chain=input ...
		detail, err := c.RunArgs(GenerateTikSentence(ruleLocation+"/print", "?", false, in))
		if err != nil {
			return ids, "", err
		}
		// get and return ids of matching rules
		for _, re := range detail.Re {
			ids = append(ids, re.Map[".id"])
		}

		return ids, ruleLocation, nil
	} else {
		return ids, "", errors.New("getIDS - missing RouterLocation tag in input interface")
	}
}

// CommitRule will execute a rule using any struct type which includes a RouterLocation
// member (IPv4FilterRule, IPv4NatRule, etc.)  Call the function with a pointer to the
// struct: c.CommitRule(&someRule).
// The struct member PlaceBeforePosition can be used to place the rule in one of these locations:
//    return  --> Place before any final return in the chain.  There can only be one action=return
//                in the chain for this to work.
//    top     --> Place at the top of the chain.
//    "comment" --> Place before the rule with the specified comment. There can only be one rule with that comment.
// If PlaceBeforePosition is blank/empty, then the rule will be placed at the bottom of the chain or before
// the rule ID specified by the PlaceBefore member.  Thus to put a rule before a specific ID, set PlaceBeforePosition
// to empty string and PlaceBefore to the ID of the rule you want to place before.
func (c *Client) CommitRule(in interface{}) error {
	// check for matching rule
	ids, location, err := c.GetIDS(in)
	if err != nil {
		return err
	}

	switch len(ids) {
	case 0:
		// get interface reflection information
		s := reflect.ValueOf(in).Elem()
		t := s.Type()

		// iterate through all fields to find position
		for i := 0; i < s.NumField(); i++ {
			field := t.Field(i)
			input := s.Field(i).Interface()

			if field.Name == "PlaceBeforePosition" {
				if len(input.(string)) > 0 {
					// place before field holder for modification
					placeBefore := s.FieldByName("PlaceBefore")

					// get id to place before if a position is specified
					switch {
					case input.(string) == "return":
						// get the chain field by name
						chain, ok := t.FieldByName("Chain")
						if !ok {
							return errors.New("commitRule - no chain field found")
						}

						// get id of the return rule
						id := []string{
							location + "/print",
							"=.proplist=.id",
							"?chain=" + s.Field(chain.Index[0]).Interface().(string),
							"?action=return",
						}

						detail, err := c.RunArgs(id)
						if err != nil {
							return err
						}

						// add place before id if a matching comment is found
						if len(detail.Re) == 1 {
							placeBefore.SetString(detail.Re[0].Map[".id"])
						} else {
							return errors.New("commitRule - more than one rule returned")
						}
					case input.(string) == "top":
						// place the rule at the top of the chain
						detail, err := c.Run(location+"/print", "=.proplist=.id")
						if err != nil {
							return err
						}

						// place before first returned rule
						placeBefore.SetString(detail.Re[0].Map[".id"])
					case len(input.(string)) > 0:
						// assume this is a comment to search for
						// get the id of the specified rule based on comment string match
						id := []string{
							location + "/print",
							"=.proplist=.id",
							"?comment=" + input.(string),
						}

						detail, err := c.RunArgs(id)
						if err != nil {
							return err
						}

						// add place before id if a matching comment is found
						if len(detail.Re) == 1 {
							placeBefore.SetString(detail.Re[0].Map[".id"])
						} else {
							return errors.New("commitRule - more than one rule returned")
						}
					}
				}
			}
		}

		// add rule to router
		_, err = c.Run(GenerateTikSentence(location+"/add", "=", true, in)...)
		if err != nil {
			return err
		}
	case 1:
		return errors.New("commitRule - matching rule exists")
	default:
		return errors.New("commitRule - more than one matching rule exists")
	}

	// return nil if all good
	return nil
}

// RemoveRule will remove a rule using any struct type which includes a RouterLocation
// member (IPv4FilterRule, IPv4NatRule, etc.)  Call the function with a pointer to the
// struct: c.RemoveRule(&someRule) to remove the rule which exactly matches that given.
func (c *Client) RemoveRule(in interface{}) error {
	// check for matching rule
	ids, location, err := c.GetIDS(in)
	if err != nil {
		return err
	}

	switch len(ids) {
	case 0:
		return errors.New("removeRule - no matching rule")
	default:
		// remove all matching rules from router
		for i := range ids {
			_, err = c.Run(location+"/remove", "=.id="+ids[i])
			if err != nil {
				return err
			}
		}
	}

	// return nil if all good
	return nil
}

// removes a rule by ID without checking if the item already exists
// can remove any rule listed in types
func (c *Client) RemoveRuleByID(in interface{}) error {
	var (
		location string
		id       string
	)
	// get interface reflection information
	s := reflect.ValueOf(in).Elem()
	t := s.Type()

	// iterate through all fields to find router location tag and the ID field
	for i := 0; i < s.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tik")
		switch {
		case field.Name == "RouterLocation":
			location = tag
		case tag == ".id":
			id = s.Field(i).Interface().(string)
		}
	}
	if len(location) > 0 && len(id) > 0 {
		_, err := c.Run(location+"/remove", "=.id="+id)
		if err != nil {
			return err
		}
	}
	return nil
}

// modify a rule using the router conn
// enable or disable rule by id
func (c *Client) ModifyRule(in interface{}, action string) error {
	// check for matching rule
	ids, location, err := c.GetIDS(in)
	if err != nil {
		return err
	}

	switch len(ids) {
	case 0:
		return errors.New("modifyRule - no matching rule")
	default:
		switch action {
		case "enable":
			// enable all matching rules from router
			for i := range ids {
				_, err = c.Run(location+"/enable", "=.id="+ids[i])
				if err != nil {
					return err
				}
			}
		case "disable":
			// disable all matching rules from router
			for i := range ids {
				_, err = c.Run(location+"/disable", "=.id="+ids[i])
				if err != nil {
					return err
				}
			}
		default:
			return errors.New("invalid modifyRule action")
		}
	}

	// return nil if all good
	return nil
}

// check if a rule is disabled
func (c *Client) RuleIsDisabled(in interface{}) (bool, error) {
	// check for matching rule
	ids, location, err := c.GetIDS(in)
	if err != nil {
		return false, err
	}

	switch len(ids) {
	case 0:
		return false, errors.New("no matching rule")
	case 1:
		detail, err := c.Run(location+"/print", "=.proplist=disabled", "?.id="+ids[0])
		if err != nil {
			return false, err
		}

		disabled, err := strconv.ParseBool(detail.Re[0].Map["disabled"])
		if err != nil {
			return false, err
		}

		// return rule state if all good
		return disabled, nil
	default:
		return false, errors.New("more than one matching rule")
	}
}

/*
// ResetConnections resets all or some connections in the connection table
func (c *Client) ResetConnections(filter string, argument string) error {
	// get all connections
	detail, err := c.Run("/ip/firewall/connection/print", "=.proplist=.id,src-address")
	if err != nil {
		log.Print(err)
	}

	// parse argument network
	_, addressNet, err := net.ParseCIDR(argument)
	if err != nil {
		log.Print(err)
	}

	// iterate through routes and match against argument
	for i := range detail.Re {
		// separate address from port
		src := strings.Split(detail.Re[i].Map["src-address"], ":")

		// remove route if src falls in argument network
		if addressNet.Contains(net.ParseIP(src[0])) {
			_, err := c.Run("/ip/firewall/connection/remove", "=.id="+detail.Re[i].Map[".id"])
			if err != nil {
				log.Print(err)
			}
		}
	}

	// return nil if all good
	return nil
}
*/

func (c *Client) AddObject(in interface{}) error {
	location := ""
	s := reflect.ValueOf(in).Elem()
	t := s.Type()

	// iterate through all fields to find position
	for i := 0; i < s.NumField(); i++ {
		field := t.Field(i)
		if field.Name == "RouterLocation" {
			location = field.Tag.Get("tik")
		}
	}
	if len(location) == 0 {
		return fmt.Errorf("no RouterLocation field in structure")
	}
	_, err := c.Run(GenerateTikSentence(location+"/add", "=", false, in)...)
	if err != nil {
		return err
	}
	return nil
}

// generate a sentence to pass to the tik api
// used to match defined rules with rules existing on the tik or set new rules
// returns a slice of strings built from the passed rule - the output can be passed directly to Run
func GenerateTikSentence(selector string, operator string, includePosition bool, i interface{}) []string {
	sentence := []string{selector}

	// get interface reflection information
	s := reflect.ValueOf(i).Elem()
	t := s.Type()

	// iterate through all fields
	for i := 0; i < s.NumField(); i++ {
		field := t.Field(i)
		tag := field.Tag.Get("tik")
		input := s.Field(i).Interface()

		// convert struct typed into strings and form a tik api compatible sentence
		// only convert if there is a tik tag
		if len(tag) > 0 {
			// skip place before field if not desired
			if !includePosition && field.Name == "PlaceBefore" {
				continue
			}

			switch field.Type.Name() {
			case "bool":
				if input.(bool) {
					sentence = append(sentence, operator+tag+"=true")
				}
			case "string":
				if len(input.(string)) > 0 {
					sentence = append(sentence, operator+tag+"="+input.(string))
				}
			case "int":
				if input.(int) > 0 {
					sentence = append(sentence, operator+tag+"="+strconv.Itoa(input.(int)))
				}
			case "time.Duration":
				d := input.(time.Duration)
				if d != 0 {
					sentence = append(sentence, operator+tag+"="+d.Round(time.Second).String())
				}
			case "time.Time":
				tm := input.(time.Time).Format("Jan/02/2006 15:04:05")
				sentence = append(sentence, operator+tag+"="+tm)
			}
		}
	}
	return sentence
}
