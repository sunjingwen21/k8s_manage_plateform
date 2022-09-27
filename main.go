package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/service"

	"github.com/gin-gonic/gin"
)

/**
 * @author 王子龙
 * 时间：2022/9/21 11:55
 */
func main() {
	//初始化gin对象
	r := gin.Default()
	//初始化k8s client
	service.K8s.Init()
	//初始化路由规则
	controller.Router.InitApiRouter(r)
	//gin程序启动
	r.Run(config.ListenAddr)
}
