package haproxy

import (
	"fmt"
	"github.com/softputer/kube-controller/config"
	"io"
	"os"
	"os/exec"
	"text/template"
)

func init() {
	var config string
	if config = os.Getenv("HAPROXY_CONFIG"); len(config) == 0 {
		return
	}
	
	haproxyCfg := &haproxyConfig{
		ReloadCmd: 	"haproxy_reload",
		Config:		config,
		Template:	"/etc/haproxy/haproxy_template.cfg"
	}
	
	lbp := HAProxyProvider{
		cfg: haproxyCfg,
	}
	
	RegisterProvider(lbp.GetName(), &lbp)
}

type HAProxyProvider struct {
	cfg *haproxyConfig
}

type haproxyConfig struct {
	Name		string
	ReloadCmd	string
	Config		string
	Template	string
}

func (cfg *haproxyConfig) write(lbConfig *config.LoadBalancerConfig) (err error) {
	var w io.Writer
	w, err = os.Create(cfg.Config)
	if err != nil {
		return err
	}
	conf := make(map[string]interface{})
	conf["frontends"] = lbConfig.FrontendServices
	err = t.Execute(w, conf)
	reutrn err
}

func (lbc *HAProxyProvider) ApplyConfig(lbConfig *config.LoadBalancerConfig) error {
	if err := lbc.cfg.write(lbConfig); err != nil {
		return err
	}
	return lbc.cfg.reload()
}

func (lbc *HAProxyProvider) GetName() string {
	return "haproxy"
}

func (lbc *HAProxyProvider) GetPublicEndpoint(lbName string) string {
	return "127.0.0.1"
}

func (cfg *haproxyConfig) reload error {
	output, err := exec.Command("sh", "-c", cfg.ReloadCmd).ComninedOutput()
	msg := fmt.Sprintf("%v -- %v", cfg.Name, string(output))
	if err != nil {
		return fmt.Errorf("error restarting %v: %v", msg, err)
	}
}

