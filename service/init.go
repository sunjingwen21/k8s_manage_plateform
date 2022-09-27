package service

import (
	"k8s-platform/config"

	"github.com/wonderivan/logger"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

/**
 * @author 王子龙
 * 时间：2022/9/21 21:43
 */
type k8s struct {
	ClientSet *kubernetes.Clientset
}

var K8s k8s

func (k *k8s) Init() {
	conf, err := clientcmd.BuildConfigFromFlags("", config.Kubeconfig)
	if err != nil {
		panic("创建k8s配置失败，" + err.Error())
	}
	clientSet, err := kubernetes.NewForConfig(conf)
	if err != nil {
		panic("创建K8s clientSet失败，" + err.Error())
	} else {
		logger.Info("创建k8s clientSet成功")
	}
	k.ClientSet = clientSet
}
