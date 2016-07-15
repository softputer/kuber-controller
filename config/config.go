package config

const (
	HTTPProto = "http"
	HTTPSProto = "https"
	TCPProto = "tcp"
	UDPProto = "udp"
)

type LoadBalancerConfig struct {
	Name			string
	Namespace		string
	FrontendServices	FrontendService
}

type FrontendService struct {
	Name			string
	Port			int
	BackendServices		BackendService
	Protocol		string
}

type BackendService struct {
	Namespace		string
	Name			string
	Port			int
	IP			string
}
