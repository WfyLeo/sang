package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type server struct {
	addr   string
	router *Router
}

func NewServer(addr string) *server {
	return &server{
		addr: addr,
	}
}

func (s *server) Router(router *Router) {
	s.router = router
}

// 启动服务器
func (s *server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http升级websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	// 允许所有CORS跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *server) wsHandler(w http.ResponseWriter, r *http.Request) {
	//websocket
	//1.http协议，省纪委websocket协议
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Println("websocket连接出错:", err)
	}
	log.Println("websocket连接成功")
	//websocket通道建立之后，不管是客户端还是服务端都可以收发消息
	//发消息的时候，把消息当做路由来处理 消息是有格式的
	//定义消息的格式
	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
	wsServer.Handshake()

}
