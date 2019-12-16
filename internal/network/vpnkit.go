package network

import "net"

// Vpnkit is an instance of Moby's vpnkit
type Vpnkit struct {
	LogDestination       string
	AllowedBindAddresses string
	Branch               string
	DB                   string
	Debug                bool
	DhcpPath             string
	Diagnostics          string
	DNS                  string
	Domain               string
	Ethernet             string
	GatewayForwards      []GatewayForward
	GatewayIP            net.IP
	GatewayNames         []string
	GcCompactInterval    int
	HighestIP            net.IP
	HostIP               net.IP
	HostNames            []string
	Hosts                string
	HTTP                 HTTPProxyConfig
	Introspection        string
	ListenBacklog        int
	LowestIP             net.IP
	MaxConnections       int
	MTU                  int
	Pcap                 string
	PcapSnaplen          int
	Port                 string
	PortMaxIdleTime      int
	ServerMacaddr        string
	TCPv4Forwards        []string
	UDPv4Forwards        []string
	VMNames              []string
	VSockPath            string
}

type DHCPConfig struct {
	/*
		{
			"searchDomains": null,
			"domainName": ""
		}
	*/
	SearchDomains []string
	DomainName    string
}

// HTTPProxyConfig is passed as a json file via the --http flag to vpnkit
type HTTPProxyConfig struct {
	TransparentHTTPPorts  []int
	TransparentHTTPSPorts []int
	Exclude               []string
}

type GatewayForward struct {
	//[{"protocol":"udp","external_port":53,"internal_ip":"127.0.0.1","internal_port":61278},{"protocol":"tcp","external_port":53,"internal_ip":"127.0.0.1","internal_port":50607}]
	Protocol     string
	ExternalPort int
	InternalIP   net.IP
	InternalPort int
}
