package controller

import (
	"fmt"
	"github.com/golang/glog"
	exampleClient "github.com/yarntime/crd-example/pkg/client/example-client"
	"github.com/yarntime/crd-example/pkg/tools"
	"github.com/yarntime/crd-example/pkg/types"
	"k8s.io/apimachinery/pkg/api/errors"
	meta_v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/wait"
	k8s "k8s.io/client-go/kubernetes"
	v1core "k8s.io/client-go/kubernetes/typed/core/v1"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/record"
	"k8s.io/client-go/util/workqueue"
	"time"
)

type ExampleController struct {
	k8sClient *k8s.Clientset

	exampleClient *exampleClient.ExampleClient

	recorder record.EventRecorder

	concurrentJobHandlers int

	resyncPeriod time.Duration

	queue workqueue.RateLimitingInterface
}

func NewExampleController(c *Config) *ExampleController {

	exampleController := &ExampleController{
		k8sClient:             c.K8sClient,
		exampleClient:         c.ExampleClient,
		concurrentJobHandlers: c.ConcurrentJobHandlers,
		resyncPeriod:          c.ResyncPeriod,
		queue:                 workqueue.NewNamedRateLimitingQueue(workqueue.DefaultControllerRateLimiter(), "example"),
	}

	eventBroadcaster := record.NewBroadcaster()
	eventBroadcaster.StartLogging(glog.Infof)
	eventBroadcaster.StartRecordingToSink(&v1core.EventSinkImpl{Interface: v1core.New(exampleController.k8sClient.CoreV1().RESTClient()).Events("")})

	_, mtlw := cache.NewInformer(
		cache.NewListWatchFromClient(exampleController.exampleClient.RESTClient(), "examples", meta_v1.NamespaceAll, fields.Everything()),
		&types.Example{},
		exampleController.resyncPeriod,
		cache.ResourceEventHandlerFuncs{
			AddFunc: exampleController.enqueueController,
			UpdateFunc: func(old, cur interface{}) {
				mt := cur.(*types.Example)
				exampleController.enqueueController(mt)
			},
		},
	)

	go mtlw.Run(c.StopCh)

	return exampleController
}

func (ec *ExampleController) Run(stopCh chan struct{}) {
	for i := 0; i < ec.concurrentJobHandlers; i++ {
		go wait.Until(ec.startHandler, time.Second, stopCh)
	}

	<-stopCh
}

func (ec *ExampleController) enqueueController(obj interface{}) {
	mt := obj.(*types.Example)
	key := tools.GetKeyOfResource(mt.ObjectMeta)
	ec.queue.Add(key)
}

func (ec *ExampleController) startHandler() {
	for ec.processNextWorkItem() {
	}
}

func (ec *ExampleController) processNextWorkItem() bool {
	key, quit := ec.queue.Get()
	if quit {
		return false
	}
	defer ec.queue.Done(key)

	ec.processExample(key.(string))
	return true
}

func (ec *ExampleController) processExample(key string) error {
	startTime := time.Now()
	defer func() {
		glog.V(4).Infof("Finished syncing example %q (%v)", key, time.Now().Sub(startTime))
	}()

	ns, name, err := cache.SplitMetaNamespaceKey(key)
	if err != nil {
		return err
	}
	if len(ns) == 0 || len(name) == 0 {
		return fmt.Errorf("invalid example key %q: either namespace or name is missing", key)
	}

	example, err := ec.exampleClient.Examples(ns).Get(name, meta_v1.GetOptions{})
	if err != nil {
		glog.Warningf("Failed get example %q (%v) from kubernetes", key, time.Now().Sub(startTime))
		if errors.IsNotFound(err) {
			glog.V(4).Infof("Example has been deleted: %v", key)
			return nil
		}
		return err
	}

	fmt.Printf("processing example:%s/%s", example.Namespace, example.Name)
	return nil
}
