package controller

import (
	"github.com/Sirupsen/logrus"
	"k8s.io/kubernetes/pkg/client/cache"
	"k8s.io/kubernetes/pkg/controller/framework"
	"k8s.io/kubernetes/pkg/util/wait"
	"k8s.io/kubernetes/pkg/util/workqueue"
	"time"
)

type TaskQueue struct {
	queue *workqueue.Type
	sync func(string)
	workerDone chan struct{}
}

func (t *TaskQueue) Run(period timeDuration, stopCh <-chan struct{}) {
	wait.Until(t.worker, period, stopCh)
}

func (t *TaskQueue) Enqueue(obj interface{}) {
	if key, ok := obj.(string); ok {
		t.queue.Add(key)
	} else {
		key, err := keyFunc(obj)
		if err != nil {
			logrus.Infof("could not get key for object %+v: %v", obj, err)
			return
		} 
		t.queue.Add(key)
	}
}

func (t *TaskQueue) Requeue(key string, err error) {
	logrus.Debug("requeuing %v, err %v", key, err)
	t.queue.Add(key)
}

func (t *TaskQueue) ShutDown() {
	t.queue.ShutDown()
	<-t.workerDone
}

func NewTaskQueue(syncFn func(string)) *TaskQueue {
	return &TaskQueue{
		queue: 		workqueue.New(),
		sync:  		syncFn,
		workerDone: 	make(chan struct{})
	}
}

