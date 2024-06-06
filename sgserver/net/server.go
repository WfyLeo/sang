package net

import (
	"github.com/gorilla/websocket"
	"log"
	"net/http"
)

type Server struct {
	addr   string  //地址
	router *router //路由
}

func NewServer(addr string) *Server {
	return &Server{
		addr: addr,
	}
}

// 启动服务
func (s *Server) Start() {
	http.HandleFunc("/", s.wsHandler)
	err := http.ListenAndServe(s.addr, nil)
	if err != nil {
		panic(err)
	}
}

// http升级为websocket协议的配置
var wsUpgrader = websocket.Upgrader{
	//允许所有cors跨域请求
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func (s *Server) wsHandler(w http.ResponseWriter, r *http.Request) {
	//websocket
	//1.http协议升级为websocket
	wsConn, err := wsUpgrader.Upgrade(w, r, nil)

	if err != nil {
		log.Fatal("websocket服务链接出错", err)
	}
	//websocket通道建立成功之后，不管是客户端还是服务端都可以接收消息
	//发消息的时候，要把消息当做路由来处理，消息是有格式的，先定义消息的格式
	//客户端发消息的时候，{name:"account.login"} 收到之后，进行解析,认为想要处理登录逻辑
	//err = wsConn.WriteMessage(websocket.BinaryMessage, []byte("hello")).
	wsServer := NewWsServer(wsConn)
	wsServer.Router(s.router)
	wsServer.Start()
}
