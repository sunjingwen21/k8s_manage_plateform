package controller

/**
 * @author 王子龙
 * 时间：2022/9/21 16:08
 */
import (
	"net/http"

	"github.com/gin-gonic/gin"
)

//实例化router结构体，可使用该对象点出首字母大写的方法（包外调用）
var Router router

//创建router结构体
type router struct{}

//初始化路由规则
func (r *router) InitApiRouter(router *gin.Engine) {
	//创建测试api接口
	router.GET("/testapi", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"msg":  "testapi success!",
			"data": nil,
		})
	}).
		////登录
		//POST("/api/login",Login.Auth).
		////工作流
		//POST("/api/k8s/workflow/create",Workflow.Create).
		//pod操作
		GET("/api/k8s/pods", Pod.GetPods).
		GET("/api/k8s/pod/detail", Pod.GetPodDetail).
		DELETE("/api/k8s/pod/del", Pod.DeletePod).
		PUT("/api/k8s/pod/update", Pod.UpdatePod).
		GET("/api/k8s/pod/container", Pod.GetPodContainer).
		GET("/api/k8s/pod/log", Pod.GetPodLog).
		GET("/api/k8s/pod/numnp", Pod.GetPodNumPerNp).
		//deployment操作
		GET("/api/k8s/deployments", Deployment.GetDeployments).
		GET("/api/k8s/deployment/detail", Deployment.GetDeploymentDetail).
		PUT("/api/k8s/deployment/scale", Deployment.ScaleDeployment).
		DELETE("/api/k8s/deployment/del", Deployment.DeleteDeployment).
		PUT("/api/k8s/deployment/restart", Deployment.RestartDeployment).
		PUT("/api/k8s/deployment/update", Deployment.UpdateDeployment).
		GET("/api/k8s/deployment/numnp", Deployment.GetDeployNumPerNp).
		POST("/api/k8s/deployment/create", Deployment.CreateDeployment)
}
