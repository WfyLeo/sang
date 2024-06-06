package main

import (
	"sg_server/config"
	"sg_server/net"
)

//http://localhost:8080/api/login

//websocket 区别 ws://localhost:8080 服务器，发消息，封装为路由

func main() {
	host := config.File.MustValue("login_server", "host", "127.0.0.1")
	port := config.File.MustValue("login_server", "port", "8003")

	s := net.NewServer(host + ":" + port)
	s.Start()
}
