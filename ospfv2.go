package gotik

import (
	"net"
	"regexp"
	"strconv"
)

type OSPF2LSA struct {
	ID         string `json:"id"`
	Instance   string `json:"instance"`
	Area       string `json:"area"`
	LSAType    string `json:"lsatype"`
	LSAID      string `json:"lsaid"`
	Originator string `json:"originator"`
	SeqNum     int    `json:"sequence-number"`
	Age        int    `json:"age"`
	Checksum   int    `json:"checksum"`
	Options    string `json:"options"`
	Body       string `json:"body"`
	Data       interface{}
}

type LSANetwork struct {
	Mask     int
	RouterID []string
}

type LSARouterLink struct {
	Type   string
	Id     string
	Data   string
	Metric int
}

type LSARouter struct {
	Flags string
	Links []LSARouterLink
}

type LSASummaryNetwork struct {
	Mask   int
	Metric int
}

type LSASummaryASBR struct {
	Metric int
}

type LSAAsExternal struct {
	Mask int
}

func parsev2lsa(props map[string]string) OSPF2LSA {
	entry := OSPF2LSA{
		ID:         props[".id"],
		Instance:   props["instance"],
		Area:       props["area"],
		LSAType:    props["type"],
		LSAID:      props["id"],
		Originator: props["originator"],
		SeqNum:     parseHex(props["sequence-number"]),
		Age:        parseInt(props["age"]),
		Checksum:   parseHex(props["checksum"]),
		Options:    props["options"],
		Body:       props["body"],
	}
	switch entry.LSAType {
	case "network":
		var x LSANetwork
		re1 := regexp.MustCompile(`netmask=([0-9.]+)\s((?:routerId=[0-9.]+\s?)+)`)
		re2 := regexp.MustCompile(`routerId=([0-9.]+)\s?`)
		if m1 := re1.FindStringSubmatch(entry.Body); m1 != nil {
			x.Mask = parseNetworkMaskToBits(m1[1])
			if m2 := re2.FindAllStringSubmatch(m1[2], -1); m2 != nil {
				for _, a := range m2 {
					x.RouterID = append(x.RouterID, a[1])
				}
			}
		}
		entry.Data = x
	case "summary-network":
		var x LSASummaryNetwork
		re := regexp.MustCompile(`netmask=([0-9.]+)\smetric=(\d+)`)
		if m1 := re.FindStringSubmatch(entry.Body); m1 != nil {
			x.Mask = parseNetworkMaskToBits(m1[1])
			x.Metric, _ = strconv.Atoi(m1[2])
		}
		entry.Data = x
	case "summary-asbr":
		var x LSASummaryASBR
		re := regexp.MustCompile(`metric=(\d+)`)
		if m1 := re.FindStringSubmatch(entry.Body); m1 != nil {
			x.Metric, _ = strconv.Atoi(m1[1])
		}
		entry.Data = x
	case "as-external":
		var x LSAAsExternal
		re := regexp.MustCompile(`netmask=([0-9.]+)`)
		if m1 := re.FindStringSubmatch(entry.Body); m1 != nil {
			x.Mask = parseNetworkMaskToBits(m1[1])
		}
		entry.Data = x
	case "router":
	}
	return entry
}

func parseNetworkMaskToBits(s string) int {
	i := net.ParseIP(s)
	m := net.IPv4Mask(i[12], i[13], i[14], i[15])
	ones, _ := m.Size()
	return ones
}

// GetOspf2LsaTable returns a slice of LSA entries on a router.  The router must be participating in OSPF
func (c *Client) GetOspf2LsaTable() ([]OSPF2LSA, error) {
	lsas := make([]OSPF2LSA, 0, 1024)
	detail, err := c.RunCmd("/routing/ospf/lsa/print")
	if err != nil {
		return lsas, err
	}
	for _, re := range detail.Re {
		lsas = append(lsas, parsev2lsa(re.Map))
	}
	return lsas, nil
}
