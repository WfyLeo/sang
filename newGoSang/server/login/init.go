package login

import (
	"newGoSang/db"
	"newGoSang/net"
	"newGoSang/server/login/controller"
)

var Router = net.NewRouter()

func Init() {
	//测试数据库，并初始化
	db.TestDB()
	//还有别的初始化方法
	initRouter()
}

func initRouter() {
	controller.DefaultAccount.Router(Router)
}
