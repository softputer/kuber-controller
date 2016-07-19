package config

const (
	HTTPProto  = "http"
	HTTPSProto = "https"
	TCPProto   = "tcp"
	UDPProto   = "udp"
)

type LoadBalancerConfig struct {
	Namespace        string
	FrontendService *FrontendService
}

type FrontendService struct {
	Name            string
	Port            int
	BackendService *BackendService
	Protocol        string
}

type BackendService struct {
	Name      string
	Port      int
	IP        string
}
