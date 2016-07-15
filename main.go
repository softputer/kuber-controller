package main

import (
	"flag"
	"github.com/softputer/kube-controller/controller"
	"github.com/softputer/kube-controller/provider"
	"os"
	"os/signal"
	"syscall"	
)

var (
	lbControllerName = flag.String("lb-controller", "kubernetes", "Ingress controller name")
	lbProviderName = flag.String("lb-provider", "haproxy", "Lb Controller name")

	lbc controller.LBController
	lbp controller.LBProvider
)


func setEnv() {
	flag.Parse()
	lbc = controller.GetController(*lbControllerName)
	if lbc == nil {
		logrus.Fatal("Unable to find controller by name %s", *lbControllerName)
	}
	lbp = provider.GetProvider(*lbProviderName)
	if lbp == nil {
		logrus.Fatal("Unable to find provider by name %s", *lbProviderName)	
	}
}

func main() {
	
}


