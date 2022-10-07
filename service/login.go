package service

import (
	"errors"
	"k8s-platform/config"

	"github.com/wonderivan/logger"
)

var Login login

type login struct{}

//验证账号密码
func (l *login) Auth(username, password string) (err error) {
	if username == config.AdminUser && password == config.AdminPwd {
		return nil
	} else {
		logger.Error("登录失败, 用户名或密码错误")
		return errors.New("登录失败, 用户名或密码错误")
	}
}
