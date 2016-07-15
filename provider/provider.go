package provider

import (
	"fmt"
	"github.com/softputer/kube-controller/config"
	 utils "github.com/softputer/kube-controller/utils"
)

const Localhost = "localhost"

type LBProvider interface {
	ApplyConfig(lbConfig *config.LoadBalancerConfig) error
	GetName() string
	GetPublicEndpoints(configName string) []string
	CleanupConfig(configName string) error
	Run(syncEndpointsQueue *util.TaskQueue)
	Stop() error
	IsHealthy() bool
}

var (
	providers map[string]LBProvider
)

func GetProvider(name string) LBProvider {
	if provider, ok := providers[name]; ok {
		return provider
	}
	return providers["haproxy"]
}

func RegisterProvider(name string, provider LBProvider) error {
	if providers == nil {
		providers = make(map[string]LBProvider)
	}
	if _, exits := providers[name]; exits {
		return fmt.Errorf("providers already registered")
	}
	providers[name] = provider
	return nil
}
