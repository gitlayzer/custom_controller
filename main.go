package main

import (
	"flag"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/fields"
	"k8s.io/apimachinery/pkg/util/runtime"
	"k8s.io/apimachinery/pkg/util/wait"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
	"k8s.io/client-go/util/workqueue"
	"k8s.io/klog/v2"
	"path/filepath"
	"time"
)

// Controller 定义一个控制器的结构体
type Controller struct {
	indexer  cache.Indexer                   // 缓存
	queue    workqueue.RateLimitingInterface // 队列
	informer cache.Controller                // 控制器
}

// NewController 用于创建一个控制器的实例
func NewController(indexer cache.Indexer, queue workqueue.RateLimitingInterface, informer cache.Controller) *Controller {
	return &Controller{
		indexer:  indexer,
		queue:    queue,
		informer: informer,
	}
}

// Run 需要定义一个 Run 方法，用于启动控制器, coroutine 用于指定启动几个 worker, stopChan 用于指定停止控制器的 channel
func (c *Controller) Run(coroutine int, stopChan chan struct{}) {
	// 用于处理 panic
	defer runtime.HandleCrash()

	// 用于关闭队列
	defer c.queue.ShutDown()

	klog.Info("starting pod controller")

	// 在这里启动 informer，用于监听资源的变化情况
	c.informer.Run(stopChan)

	if !cache.WaitForCacheSync(stopChan, c.informer.HasSynced) {
		klog.Errorf("timed out waiting for caches to sync")
		return
	}

	// 在这里启动 worker，可以启动多个，是以 goroutine 的方式启动
	for i := 0; i < coroutine; i++ {
		// 使用 wait.Until 方法，启动一个 goroutine，用于调用 RunWorker 方法
		go wait.Until(c.RunWorker, time.Second, stopChan)
	}

	// 等待停止信号
	<-stopChan

	klog.Info("stopping pod controller")
}

// RunWorker 用于启动一个 worker，用于处理队列中的元素
func (c *Controller) RunWorker() {
	// 一直运行，直到队列关闭
	for c.ProcessItem() {
	}
}

// ProcessItem 获取队列中的元素，然后调用 HandlerObj 方法处理元素
func (c *Controller) ProcessItem() bool {
	key, quit := c.queue.Get() // 获取队列中的元素
	// 如果队列关闭，则返回 false
	if quit {
		return false
	}

	// 处理完元素后，调用 Done 方法，以便通知队列，当前元素已经处理完毕
	defer c.queue.Done(key)

	// 处理元素
	if err := c.HandlerObject(key.(string)); err != nil {
		// 如果处理失败次数小于 5 次，则重新添加到队列中
		if c.queue.NumRequeues(key) < 5 {
			c.queue.Add(key)
		}
	}

	return true
}

// HandlerObject 用于处理事件，它是通过 queue 传递过来的 key 来获取 Indexer（缓存）中的对象（Object）
func (c *Controller) HandlerObject(key string) error {
	// 通过 key 获取缓存中的对象
	obj, exists, err := c.indexer.GetByKey(key)
	if err != nil {
		klog.Errorf("fetch object from store failed: %v", err)
		return err
	}

	// 如果对象不存在，则可能是对象已经被删除了，所以这里需要删除对应的资源
	if !exists {
		// 对象不存在
		klog.Infof("object %s does not exist in store", key)
	} else {
		// 处理对象
		klog.Infof("handle object %s", obj.(*v1.Pod).GetName(), obj.(*v1.Pod).GetNamespace())
	}

	return nil
}

// initClient 获取 k8s client
func initClient() (*kubernetes.Clientset, error) {
	var err error           // 错误变量
	var config *rest.Config // 配置变量
	var kubeConfig *string  // kubeConfig 文件路径

	// 获取 kubeConfig 文件路径
	if home := homedir.HomeDir(); home != "" {
		kubeConfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeConfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse() // 解析命令行参数

	if config, err = rest.InClusterConfig(); err != nil {
		if config, err = clientcmd.BuildConfigFromFlags("", *kubeConfig); err != nil {
			return nil, err
		}
	}

	// 返回 k8s client
	return kubernetes.NewForConfig(config)
}

func main() {
	// 初始化 k8s client
	clientSet, err := initClient()
	if err != nil {
		klog.Fatal(err)
	}

	// queue 用于创建一个队列
	queue := workqueue.NewRateLimitingQueue(workqueue.DefaultControllerRateLimiter())

	// 创建一个 Pod 的缓存器
	// clientSet.CoreV1().RESTClient() 用于获取 Pod 的 RESTClient
	// "pods" 用于指定监听哪个资源
	// v1.NamespaceAll 用于指定监听所有的 namespace
	// fields.Everything() 用于指定监听所有的字段
	podListWatch := cache.NewListWatchFromClient(clientSet.CoreV1().RESTClient(), "pods", v1.NamespaceAll, fields.Everything())

	// indexer, informer 用于创建一个缓存和控制器
	// podListWatch 用于指定监听哪个资源
	// &v1.Pod{} 用于指定监听资源的类型
	// 0 用于指定 resyncPeriod，如果为 0，则表示不会主动同步，只有当有事件发生时，才会同步
	// cache.ResourceEventHandlerFuncs{} 用于指定事件处理函数
	// cache.Indexers{} 用于指定索引器
	indexer, informer := cache.NewIndexerInformer(podListWatch, &v1.Pod{}, 0, cache.ResourceEventHandlerFuncs{
		// 当有资源添加时，会调用该方法
		AddFunc: func(obj interface{}) {
			// 将资源的 key 添加到队列中
			key, err := cache.MetaNamespaceKeyFunc(obj)
			klog.Info("AddFunc", key)
			if err == nil {
				// 将资源的 key 添加到队列中
				queue.Add(key)
			}
		},
		// 当有资源更新时，会调用该方法
		UpdateFunc: func(oldObj, newObj interface{}) {
			// 如果资源的 ResourceVersion 没有变化，则不处理
			if oldObj.(*v1.Pod).ResourceVersion == newObj.(*v1.Pod).ResourceVersion {
				return
			} else {
				// 将资源的 key 添加到队列中
				key, err := cache.MetaNamespaceKeyFunc(newObj)
				klog.Info("UpdateFunc", key)
				if err == nil {
					// 将资源的 key 添加到队列中
					queue.Add(key)
				}
			}
		},
		// 当有资源删除时，会调用该方法
		DeleteFunc: func(obj interface{}) {
			// 将资源的 key 添加到队列中
			key, err := cache.DeletionHandlingMetaNamespaceKeyFunc(obj)
			klog.Info("DeleteFunc", key)
			if err == nil {
				// 将资源的 key 添加到队列中
				queue.Add(key)
			}
		},
	}, cache.Indexers{})

	// 创建一个控制器
	controller := NewController(indexer, queue, informer)

	// 创建一个 channel，用于通知控制器停止工作
	stopChan := make(chan struct{})

	// 启动控制器
	go controller.Run(1, stopChan)

	// 等待停止信号
	defer close(stopChan)
	select {}
}
