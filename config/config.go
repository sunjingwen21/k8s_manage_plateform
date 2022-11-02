package config

import "time"

/**
 * @author 王子龙
 * 时间：2022/9/21 16:11
 */
const (
	//监听端口
	ListenAddr = "0.0.0.0:9093"
	//权限文件地址
	Kubeconfig = "/root/.kube/config"
	//数据库配置
	DbType = "mysql"
	DbUser = "root"
	DbPwd  = "123456"
	DbHost = "10.0.12.12"
	DbPort = 30002
	DbName = "test_db"
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
