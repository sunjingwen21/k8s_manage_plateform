package main

import (
	"k8s-platform/config"
	"k8s-platform/controller"
	"k8s-platform/db"
	"k8s-platform/middle"
	"k8s-platform/service"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
 * @author 王子龙
 * 时间：2022/9/21 11:55
 */
func main() {
	//初始化k8s clientset
	service.K8s.Init()
	//初始化数据库
	db.Init()
	//关闭db连接
	defer db.Close()
	//初始化gin对象路由配置
	r := gin.Default()
	//跨域配置
	r.Use(middle.Cors())
	//jwt token验证
	r.Use(middle.JWTAuth())
	//初始化路由规则
	controller.Router.InitApiRouter(r)

	//终端websocket
	go func() {
		http.HandleFunc("/ws", service.Terminal.WsHandler)
		http.ListenAndServe(":9094", nil)
	}()

	//http server gin程序启动
	r.Run(config.ListenAddr)
}
