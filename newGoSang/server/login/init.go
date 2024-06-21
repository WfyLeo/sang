package login

import (
	"newGoSang/net"
	"newGoSang/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	//还有别的初始化方法
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
