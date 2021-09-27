package gotik

import (
	"time"
)

type IPv4Address struct {
	ID              string `json:"id"`
	Address         string `json:"address"` // CIDR format
	Network         string `json:"network"`
	Interface       string `json:"intf"`
	ActualInterface string `json:"actual-intf"`
	Invalid         bool   `json:"invalid"`
	Dynamic         bool   `json:"dynamic"`
	Disabled        bool   `json:"disabled"`
	Comment         string `json:"comment"`
}

type IPv4Pool struct {
	ID       string `json:"id"`
	Name     string `json:"name"`
	Ranges   string `json:"ranges"`
	NextPool string `json:"nextpool"`
}

type DHCP4Network struct {
	ID         string `json:"id"`
	Address    string `json:"address"`
	Gateway    string `json:"gateway"`
	Netmask    string `json:"netmask"`
	Domain     string `json:"domain"`
	Comment    string `json:"comment"`
	DNSServers string `json:"dnsservers"` // Address,Address
	NTPServers string `json:"ntpservers"` // Address,Address
	Options    string `json:"options"`    // Option,Option
}

type DHCPv4Server struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	LeaseTime     string `json:"leasetime"` // time interval
	Pool          string `json:"pool"`
	Interface     string `json:"intf"`
	Authoritative string `json:"authoritative"` // yes, no, after-10sec-delay, after-2sec-delay
	Disabled      bool   `json:"disabled"`
	AddArp        bool   `json:"addarp"`
}

const (
	InterfaceType6to4            = "6to4"
	InterfaceTypeBonding         = "bonding"
	InterfaceTypeBridge          = "bridge"
	InterfaceTypeEoip            = "eoip"
	InterfaceTypeEoipv6          = "eoipv6"
	InterfaceTypeEthernet        = "ether"
	InterfaceTypeGre             = "gre"
	InterfaceTypeGre6            = "gre6"
	InterfaceTypeIpip            = "ipip"
	InterfaceTypeIpipv6          = "ipipv6"
	InterfaceTypeL2tpClient      = "l2tp-client"
	InterfaceTypeL2tpServer      = "l2tp-server"
	InterfaceTypeLte             = "lte"
	InterfaceTypeMesh            = "mesh"
	InterfaceTypeOvpnClient      = "ovpn-client"
	InterfaceTypeOvpnServer      = "ovpn-server"
	InterfaceTypePppClient       = "ppp-client"
	InterfaceTypePppServer       = "ppp-server"
	InterfaceTypePppoeClient     = "pppoe-client"
	InterfaceTypePppoeServer     = "pppoe-server"
	InterfaceTypePptpClient      = "pptp-client"
	InterfaceTypePptpServer      = "pptp-server"
	InterfaceTypeSstpClient      = "sstp-client"
	InterfaceTypeSstpServer      = "sstp-server"
	InterfaceTypeTrafficEng      = "traffic-eng"
	InterfaceTypeVirtualEthernet = "virtual-ethernet"
	InterfaceTypeVlan            = "vlan"
	InterfaceTypeVpls            = "vpls"
	InterfaceTypeVrrp            = "vrrp"
	InterfaceTypeWireless        = "wireless"
	InterfaceTypeWDS             = "wds"
)

type Interface struct {
	ID            string `json:"id"`
	Type          string `json:"intftype"`
	Name          string `json:"name"`
	Mac           string `json:"mac"`
	OriginalMac   string `json:"origmac"`
	Interface     string `json:"interface"` // base interface for certain types
	Disabled      bool   `json:"disabled"`
	Dynamic       bool   `json:"dynamic"`
	Running       bool   `json:"running"`
	Arp           string `json:"arp"`
	VLAN          int    `json:"vlanid"`
	Comment       string `json:"comment"`
	AutoNeg       bool   `json:"autoneg"`     // auto-negotiate enabled?
	Speed         string `json:"speed"`       // Hard-set speed, only applicable if AutoNeg is false
	FullDuplex    bool   `json:"fullduplex"`  // Hard-set duplex, only applicable if AutoNeg is false
	DefaultName   string `json:"defaultname"` // Only applicable to Ethernet
	Slave         bool   `json:"slave"`
	AdminMac      string `json:"adminmac"`
	AutoMac       bool   `json:"automac"`
	ProtocolMode  string `json:"protocolmode"`
	AgingTime     string `json:"agingtime"`
	VlanFiltering bool   `json:"vlanfiltering"`
	FastForward   bool   `json:"fastforward"`
}

type EthernetStatus struct {
	ID             string `json:"id"`
	Name           string `json:"name"`
	Status         string `json:"status"`
	AutoNegStatus  string `json:"auto_neg_status"`
	Rate           string `json:"rate"`
	FullDuplex     bool   `json:"full_duplex"`
	TxFlowControl  bool   `json:"tx_flow_control"`
	RxFlowControl  bool   `json:"rx_flow_control"`
	Advertising    string `json:"advertising"`
	LinkPartnerAdv string `json:"link_partner_adv"`
}

type ArpEntry struct {
	ID        string `json:"id"`
	Address   string `json:"ip4addr"`
	Mac       string `json:"mac"`
	Interface string `json:"interface"`
	Comment   string `json:"comment"`
	Disabled  bool   `json:"disabled"`
	Dynamic   bool   `json:"dynamic"`
	Complete  bool   `json:"complete"`
	DHCP      bool   `json:"dhcp"`
}

type IPv4Route struct {
	ID            string `json:"id"`
	Gateway       string `json:"gateway"`
	GatewayStatus string `json:"gwstatus"`
	DstAddress    string `json:"dstaddress"`
	Comment       string `json:"comment"`
	PrefSrc       string `json:"prefsrc"`
	Mark          string `json:"mark"`
	Distance      int    `json:"distance"`
	Scope         int    `json:"scope"`
	TargetScope   int    `json:"targetscope"`
	RouteType     string `json:"routetype"` // blackhole, prohibit, unicast, unreachable
	Active        bool   `json:"active"`
	Disabled      bool   `json:"disabled"`
	Static        bool   `json:"static"`
	Connected     bool   `json:"connected"`
}

type QueueTree struct {
	ID             string      `json:"id"`
	BucketSize     float32     `json:"bucketsize"`     // 0..10 (defaults to 0.1/0.1)
	BurstLimit     int         `json:"burstlimit"`     // in bps
	BurstThreshold int         `json:"burstthreshold"` // in bps
	BurstTime      string      `json:"bursttime"`
	Comment        string      `json:"comment"`
	Disabled       bool        `json:"disabled"`
	Dynamic        bool        `json:"dynamic"`
	Invalid        bool        `json:"invalid"`
	Name           string      `json:"name"`
	LimitAt        int         `json:"limitat"`  // in bps
	MaxLimit       int         `json:"maxlimit"` // in bps
	PacketMark     string      `json:"packetmark"`
	Parent         string      `json:"parent"`
	Priority       int         `json:"priority"` // 1..8, only valid if Parent > ''
	Queue          string      `json:"queue"`    // type of queue
	Children       []QueueTree `json:"children"`
}

type SimpleQueue struct {
	ID             string     `json:"id"`
	BucketSize     [2]float32 `json:"bucketsize"`     // Upload/Download 0..10 (defaults to 0.1/0.1)
	BurstLimit     [2]int     `json:"burstlimit"`     // Upload/Download in bps
	BurstThreshold [2]int     `json:"burstthreshold"` // Upload/Download in bps
	BurstTime      [2]string  `json:"bursttime"`      // Upload/Download
	Comment        string     `json:"comment"`
	Disabled       bool       `json:"disabled"`
	Dynamic        bool       `json:"dynamic"`
	Dst            string     `json:"dst"`
	Invalid        bool       `json:"invalid"`
	LimitAt        [2]int     `json:"limitat"`  // Upload/Download in bps
	MaxLimit       [2]int     `json:"maxlimit"` // Upload/Download in bps
	Name           string     `json:"name"`
	PacketMarks    string     `json:"packetmarks"`
	Parent         string     `json:"parent"`
	Priority       int        `json:"priority"` // 1..8, only valid if Parent > ''
	Queue          [2]string  `json:"queue"`    // type of queue
	Target         string     `json:"target"`
	Time           string     `json:"time"`
}

type PPPSecret struct {
	ID            string `json:"id"`
	Name          string `json:"name"`
	CallerID      string `json:"callerid"`
	Comment       string `json:"comment"`
	Disabled      bool   `json:"disabled"`
	LimitBytesIn  int    `json:"limitbytesin"`
	LimitBytesOut int    `json:"limitbytesout"`
	LocalAddress  string `json:"localaddress"`
	Password      string `json:"password"`
	Profile       string `json:"profile"`
	RemoteAddress string `json:"remoteaddress"`
	Routes        string `json:"routes"`
	Service       string `json:"service"`
}

type PPPActive struct {
	ID            string        `json:"id"`
	Name          string        `json:"name"`
	CallerID      string        `json:"callerid"`
	SessionID     int           `json:"sessionid"`
	Address       string        `json:"address"`
	Service       string        `json:"service"`
	Radius        bool          `json:"radius"`
	Uptime        time.Duration `json:"uptime"`
	Encoding      string        `json:"encoding"`
	LimitBytesIn  int           `json:"limitbytesin"`
	LimitBytesOut int           `json:"limitbytesout"`
}

type PPPoEServer struct {
	ID             string `json:"id"`
	Disabled       bool   `json:"disabled"`
	Interface      string `json:"interface"`
	ServiceName    string `json:"service-name"`
	MaxMTU         int    `json:"max-mtu"`
	MaxMRU         int    `json:"max-mru"`
	MRRU           int    `json:"mrru"`
	Authentication string `json:"authentication"`
	KeepAlive      int    `json:"keepalive-timeout"`
	SingleSess     bool   `json:"one-session-per-host"`
	MaxSessions    string `json:"max-sessions"`
	DefaultProfile string `json:"default-profile"`
	PadoDelay      int    `json:"pado-delay"`
}

type NeighborInterface struct {
	ID       string
	Discover string
}

type AddressList struct {
	ID             string `tik:".id"`
	List           string `tik:"list"`
	Dynamic        bool
	Disabled       bool   `tik:"disabled"`
	Address        string `tik:"address"`
	Comment        string `tik:"comment"`
	CreationTime   time.Time
	Timeout        time.Duration `tik:"timeout"`
	RouterLocation string        `tik:"/ip/firewall/address-list"`
}

type IPv4NatRule struct {
	RouterLocation      string `tik:"/ip/firewall/nat"`
	ID                  string `tik:".id"`
	PlaceBeforePosition string
	PlaceBefore         string `tik:"place-before"`
	Disabled            bool   `tik:"disabled"`
	Dynamic             bool   `tik:"dynamic"`
	Invalid             bool   `tik:"invalid"`
	Chain               string `tik:"chain"`
	SrcAddress          string `tik:"src-address"`
	DstAddress          string `tik:"dst-address"`
	Protocol            string `tik:"protocol"`
	SrcPort             int    `tik:"src-port"`
	DstPort             int    `tik:"dst-port"`
	InInterface         string `tik:"in-interface"`
	OutInterface        string `tik:"out-interface"`
	SrcAddressList      string `tik:"src-address-list"`
	DstAddressList      string `tik:"dst-address-list"`
	ToAddresses         string `tik:"to-addresses"`
	ToPorts             int    `tik:"to-ports"`
	Comment             string `tik:"comment"`
	Action              string `tik:"action"`
	JumpTarget          string `tik:"jump-target"`
}

type IPv4FilterRule struct {
	RouterLocation          string `tik:"/ip/firewall/filter"`
	ID                      string `tik:".id"`
	PlaceBeforePosition     string
	Action                  string        `tik:"action"`
	AddressList             string        `tik:"address-list"`
	AddressListTimeout      time.Duration `tik:"address-list-timeout"`
	Chain                   string        `tik:"chain"`
	Comment                 string        `tik:"comment"`
	Disabled                bool          `tik:"disabled"`
	Dynamic                 bool          `tik:"dynamic"`
	Invalid                 bool          `tik:"invalid"`
	DstAddress              string        `tik:"dst-address"`
	DstAddressList          string        `tik:"dst-address-list"`
	DstPort                 string        `tik:"dst-port"`
	InInterface             string        `tik:"in-interface"`
	InInterfaceList         string        `tik:"in-interface-list"`
	JumpTarget              string        `tik:"jump-target"`
	Log                     bool          `tik:"log"`
	LogPrefix               string        `tik:"log-prefix"`
	OutInterface            string        `tik:"out-interface"`
	OutInterfaceList        string        `tik:"out-interface-list"`
	PlaceBefore             string        `tik:"place-before"`
	Protocol                string        `tik:"protocol"`
	RejectWith              string        `tik:"reject-with"`
	SrcAddress              string        `tik:"src-address"`
	SrcAddressList          string        `tik:"src-address-list"`
	SrcPort                 string        `tik:"src-port"`
	TcpFlags                string        `tik:"tcp-flags"`
	TcpMss                  string        `tik:"tcp-mss"`
	ConnectionBytes         string        `tik:"connection-bytes"`          // Match packets with given bytes or byte range
	ConnectionLimit         string        `tik:"connection-limit"`          // Restrict connection limit per address or address block
	ConnectionMark          string        `tik:"connection-mark"`           // Matches packets marked via mangle facility with particular connection mark
	ConnectionNatState      string        `tik:"connection-nat-state"`      // dstnat, srcnat, !dstnat, !srcnat
	ConnectionRate          string        `tik:"connection-rate"`           // ConnectionRate ::= [!]From,To ::= 0..4294967295
	ConnectionState         string        `tik:"connection-state"`          // Interprets the connection tracking analysis data for a particular packet
	ConnectionType          string        `tik:"connection-type"`           // Match packets with given connection type
	Content                 string        `tik:"content"`                   // The text packets should contain in order to match the rule
	DSCP                    string        `tik:"dscp"`                      //
	DstAddressType          string        `tik:"dst-address-type"`          // Destination address type
	DstLimit                string        `tik:"dst-limit"`                 // Packet limitation per time with burst to dst-address, dst-port or src-address
	Fragment                string        `tik:"fragment"`                  //
	Hotspot                 string        `tik:"hotspot"`                   // Matches packets received from clients against various Hot-Spot
	IcmpOptions             string        `tik:"icmp-options"`              // IcmpOptions ::= [!]Type[:Code]; Type ::= 0..255; Code ::= Start[-End] ::= 0..255
	InBridgePort            string        `tik:"in-bridge-port"`            //
	InBridgePortList        string        `tik:"in-bridge-port-list"`       //
	IngressPriority         string        `tik:"ingress-priority"`          // IngressPriority ::= [!]IngressPriority ::= 0..63
	IpsecPolicy             string        `tik:"ipsec-policy"`              //
	Ipv4Options             string        `tik:"ipv4-options"`              // Match ipv4 header options
	Layer7Protocol          string        `tik:"layer7-protocol"`           //
	Limit                   string        `tik:"limit"`                     // Setup burst, how many times to use it in during time interval measured in seconds
	Nth                     string        `tik:"nth"`                       // Match nth packets received by the rule
	OutBridgePort           string        `tik:"out-bridge-port"`           // Matches the bridge port physical output device added to a bridge device
	OutBridgePortList       string        `tik:"out-bridge-port-list"`      //
	PacketMark              string        `tik:"packet-mark"`               // Matches packets marked via mangle facility with particular packet mark
	PacketSize              string        `tik:"packet-size"`               // Packet size or range in bytes
	PerConnectionClassifier string        `tik:"per-connection-classifier"` //
	Port                    string        `tik:"port"`                      //
	Priority                string        `tik:"priority"`                  //
	PSD                     string        `tik:"psd"`                       // Detect TCP un UDP scans
	PktRandom               string        `tik:"random"`                    // Match packets randomly with given propability
	RoutingMark             string        `tik:"routing-mark"`              // Matches packets marked by mangle facility with particular routing mark
	RoutingTable            string        `tik:"routing-table"`             //
	SrcAddressType          string        `tik:"src-address-type"`          // Source IP address type
	SrcMacAddress           string        `tik:"src-mac-address"`           // Source MAC address
	PktTime                 string        `tik:"time"`                      // Packet arrival time and date or locally generated packets departure time and date
	TLSHost                 string        `tik:"tls-host"`                  //
	TTL                     string        `tik:"ttl"`                       //
}

type PackageUpdate struct {
	Channel   string `json:"channel"`
	Installed string `json:"installed"`
	Latest    string `json:"latest"`
	Status    string `json:"status"`
}
