package controller

import (
	"fmt"
	"github.com/softputer/kuber-controller/provider"
	"github.com/softputer/kuber-controller/config"
)

type LBController interface {
	GetName() string
	Run(lbProvider provider.LBProvider)
	Stop() error
	GetLBConfigs() []*config.LoadBalancerConfig
	IsHealthy() bool
}

var (
	controllers map[string]LBController
)

func GetController(name string) LBController {
	if controller, ok := controllers[name]; ok {
		return controller
	}
	return controllers["kubernetes"]
}

func RegisterController(name string, controller LBController) error {
	if controllers == nil {
		controllers = make(map[string]LBController)
	}
	if _, exists := controllers[name]; exits {
		return fmt.Errorf("controller already registered")
	}
	controllers[name] = controller
	return nil
}
