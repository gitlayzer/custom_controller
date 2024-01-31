# custom_controller
Understand the informer mechanism

# This project is only for learning the informer mechanism

![](https://img-blog.csdnimg.cn/5ebf710b7c804d2abc0303a999d80219.png)


# premise
```shell
You must have a kubeConfig for you to connect to the cluster
```

# Sample
```shell
❯ go run main.go                            
I0201 04:30:51.213610   16024 main.go:44] starting pod controller
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-flannel/kube-flannel-ds-qcml7  
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-flannel/kube-flannel-ds-s6xhq  
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-system/coredns-5d78c9869d-k5274
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/coredns-5d78c9869d-rwcg4
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/etcd-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-apiserver-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-controller-manager-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-proxy-fhgfw
I0201 04:30:51.237918   16024 main.go:168] AddFunckube-system/kube-proxy-gb9s5
I0201 04:30:51.237918   16024 main.go:168] AddFunckube-system/kube-scheduler-gitlayzer-control-plane
I0201 04:30:51.237918   16024 main.go:168] AddFunclocal-path-storage/local-path-provisioner-6bc4bddd6b-qlnp6

# Get all Pods and add them to the cache

# At this time, creating, updating, and deleting Pods will trigger the ResourceEventHandlerFuncs function.

# Create Pod
root@ubuntu:~# kubectl apply -f pod.yaml 
pod/nginx created

❯ go run main.go                            
I0201 04:30:51.213610   16024 main.go:44] starting pod controller
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-flannel/kube-flannel-ds-qcml7  
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-flannel/kube-flannel-ds-s6xhq  
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-system/coredns-5d78c9869d-k5274
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/coredns-5d78c9869d-rwcg4
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/etcd-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-apiserver-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-controller-manager-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-proxy-fhgfw
I0201 04:30:51.237918   16024 main.go:168] AddFunckube-system/kube-proxy-gb9s5
I0201 04:30:51.237918   16024 main.go:168] AddFunckube-system/kube-scheduler-gitlayzer-control-plane
I0201 04:30:51.237918   16024 main.go:168] AddFunclocal-path-storage/local-path-provisioner-6bc4bddd6b-qlnp6
I0201 04:32:45.503556   16024 main.go:168] AddFuncdefault/nginx     # AddFunc
I0201 04:32:45.508018   16024 main.go:182] UpdateFuncdefault/nginx  # UpdateFunc
I0201 04:32:45.520421   16024 main.go:182] UpdateFuncdefault/nginx  # UpdateFunc
I0201 04:32:46.242495   16024 main.go:182] UpdateFuncdefault/nginx  # UpdateFunc

# Delete Pod
root@ubuntu:~# kubectl delete -f pod.yaml 
pod "nginx" deleted


❯ go run main.go                            
I0201 04:30:51.213610   16024 main.go:44] starting pod controller
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-flannel/kube-flannel-ds-qcml7  
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-flannel/kube-flannel-ds-s6xhq  
I0201 04:30:51.236952   16024 main.go:168] AddFunckube-system/coredns-5d78c9869d-k5274
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/coredns-5d78c9869d-rwcg4
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/etcd-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-apiserver-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-controller-manager-gitlayzer-control-plane
I0201 04:30:51.237421   16024 main.go:168] AddFunckube-system/kube-proxy-fhgfw
I0201 04:30:51.237918   16024 main.go:168] AddFunckube-system/kube-proxy-gb9s5
I0201 04:30:51.237918   16024 main.go:168] AddFunckube-system/kube-scheduler-gitlayzer-control-plane
I0201 04:30:51.237918   16024 main.go:168] AddFunclocal-path-storage/local-path-provisioner-6bc4bddd6b-qlnp6
I0201 04:32:45.503556   16024 main.go:168] AddFuncdefault/nginx
I0201 04:32:45.508018   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:32:45.520421   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:32:46.242495   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:33:46.656692   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:33:46.862729   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:33:47.353657   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:33:47.360599   16024 main.go:182] UpdateFuncdefault/nginx
I0201 04:33:47.362585   16024 main.go:193] DeleteFuncdefault/nginx   # DeleteFunc
```
