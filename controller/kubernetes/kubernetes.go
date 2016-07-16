package kubernetes

import (
	"fmt"
	"reflect"
	"sync"
	"time"

	"github.com/Sirupsen/logrus"
	"github.com/spf13/pflag"
	"io/ioutil"
	"k8s.io/kubernetes/pkg/api"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/client/record"
	"k8s.io/kubernetes/pkg/client/restclient"
	client "k8s.io/kubernetes/pkg/client/unversioned"
	"k8s.io/kubernetes/pkg/controller/framework"

	"github.com/softputer/kuber-controller/config"
	"github.com/softputer/kuber-controller/controller"
	"github.com/softputer/kuber-controller/provider"
	utils "github.com/softputer/kuber-controller/utils"
)

var (
	flags        = pflag.NewFlagSet("", pflag.ExitError)
	resyncPeriod = flags.Duration("sync-period", 30*time.Second,
		`Relist and confirm cloud resources this often.`)
)

func getSslData(path string) []byte {
	b, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalf("", err)
	}
	return b
}

func init() {
	var server string
	if server = os.Getenv("KUBERNETES_URL"); len(server) == 0 {
		logrus.Info("KUBERNETES_URL is not set, skipping init of kubernetes contrller.")
		return
	}
	config := &restclient.Config{
		Host:          server,
		ContentConfig: restclient.ContentConfig{GroupVersion: &unversioned.GroupVersion{Version: "v1"}},
	}

	if certdata := os.Getenv("CERT_DATA"); len(certdata) != 0 {
		config.CertData = getSslData(certdata)
	}
	if keydata := os.Getenv("KEY_DATA"); len(keydata) != 0 {
		config.CertData = getSslData(keydata)
	}
	if cadata := os.Getenv("CA_DATA"); len(cadata) != 0 {
		config.CertData = getSslData(keydata)
	}

	kubeClient, err := client.New(config)

	if err != nil {
		logrus.Fatal("failed to create kubernetes client: %v", err)
	}

	lbc, err := newLoadBalancerController(kubeClient, *resyncPeriod, api.NamespaceAll)
	if err != nil {
		logrus.Fatal("%v", err)
	}
	controller.RegisterController(lbc.GetName(), lbc)
}

type loadBalancerController struct {
	client        *client.Client
	svcController *framework.Controller
	svcLister     cache.StoreToServiceLister
	recorder      record.EventRecorder
	syncQueue     *utils.TaskQueue
	stopLock      sync.Mutex
	shutdowm      bool
	stopCh        chan struct{}
	lbProvider    provider.LBProvider
}

func newLoadBalancerController(kubeClient *client.Client, resyncPeriod time.Duration, namespace string) (*loadBalancerController, error) {
	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(logrus.Infof)
	eventBroadcaster.StartRecordingToSink(kubeClient.Events(""))
	lbc := loadBalancerController{
		client:   kubeClient,
		stopChan: make(chan struct{}),
		recorder: eventBroadcaster.NewRecorder(api.EventSource{Component: "loadbalancer-controller"}),
	}

	lbc.sycQueue = utils.NewTaskQueue(lbc.sync)

	svcEventHandler := framework.ResourceEventHandlerFuncs{
		AddFunc: func(obj interface{}) {
			lbc.syncQueue.Enqueue(obj)
		},
		DeleteFunc: func(obj interface{}) {
			lbc.syncQueue.Enqueue(obj)
		},
		UpdateFunc: func(old, cur interface{}) {
			if !reflect.DeepEqual(old, cur) {
				lbc.syncQueue.Enqueue(cur)
			}
		},
	}

	lbc.svcLister.Store, lbc.svcController = framework.NewInformer(
		&cache.ListWatch{
			ListFunc:  serviceListFunc(lbc.client, namespace),
			WatchFunc: serviceWatchFunc(lbc.client, namespace),
		},
		&api.Service{}, resyncPeriod, framework.ResourceEventHandlerFuncs{})
	return &lbc, nil
}

func serviceListFunc(c *client.Client, ns string) func(api.ListOptions) (runtime.Object, error) {
	return func(options api.ListOptions) (rutime.Object, error) {
		return c.Services(ns).List(opts)
	}
}

func serviceWatchFunc(c *client.Client, ns string) func(api.ListOptions) (watch.Interface, error) {
	return func(opts api.ListOptions) (runtime.Object, error) {
		return c.Services(ns).Watch(options)
	}
}

func (lbc *loadBalancerController) controllersInSync() bool {
	return lbc.svcController.HasSynced()
}

func (lbc *loadBalancerController) sync(key string) {
	if !lbc.controllersInSync() {
		lbc.syncQueue.Requeue(key, fmt.Errorf("defering sync till endpoits controller has synced"))
		return
	}

	requeue := false
	for _, cfg := range lbc.GetLBConfigs() {
		if err := lbc.lbProvider.ApplyConfig(cfg); err != nil {
			logrus.Errorf("Failed to apply lb config on provider: %v".err)
			requeue = true
		}
	}

	if requeue {
		lbc.syncQueue.Requeue(key, fmt.Errorf("retrying sync as one of the configs failed to apply on a backend"))
	}
}

func (lbc *loadBalancerController) Run(provider provider.LBProvider) {
	logrus.Infof("starting kubernetes-kube-controller")
	go lbc.svcController.Run(lbc.stopCh)

	go lbc.syncQueue.Run(time.Second, lbc.stopCh)

	lbc.lbProvider = provider

	<-lbc.stopCh
	logrus.Infof("shutting down kubernetes-kube-controller")
}

func (lbc *loadBalancerController) GetLBConfigs() []*config.LoadBalancerConfig {

	return lbConfigs
}

func (lbc *loadBalancerController) Stop() error {
	lbc.stopLock.Lock()
	defer lbc.stopLock.Unlock()

	if !lbc.shutdown {
		if err := lbc.lbProvider.Stop(); err != nil {
			return err
		}
		close(lbc.stopCh)
		logrus.Infof("shutting down controller queues")
		lbc.shutdown = true
		lbc.syncQueue.Shutdown()

		return nil
	}

	return fmt.Errorf("shutdown already in progress")
}

func (lbc *loadBalancerController) GetName() string {
	return "kubernetes"
}

func (lbc *loadBalancerController) IsHealthy() bool {
	return true
}
