package config

import "time"

const (
	//监听端口
	ListenAddr = "0.0.0.0:9093"
	//权限文件地址
	//if os = Windows
	Kubeconfig = "C:\\Users\\13358\\.kube\\config"
	//else if os = linux
	//Kubeconfig = "/root/.kube/config"
	//数据库配置
	DbType = "mysql"
	DbUser = "root"
	DbPwd  = "123456"
	DbHost = "127.0.0.1"
	DbPort = 31157
	DbName = "platform"
	//连接池的配置
	MaxIdleConns = 10               //最大空闲连接
	MaxOpenConns = 100              //最大连接数
	MaxLifeTime  = 30 * time.Second //最大生存时间
	//日志显示行数
	PodLogTailLine = 2000
	//登录账户名和密码
	AdminUser = "admin"
	AdminPwd  = "qwer1234"
)
