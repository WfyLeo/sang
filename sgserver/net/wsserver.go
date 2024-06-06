package net

import (
	"fmt"
	"github.com/gorilla/websocket"
	"log"
	"sync"
)

// websocket服务
type wsServer struct {
	wsConn       *websocket.Conn //webSocket链接
	router       *router         //路由
	outChan      chan *WsMsgRsp  //通道
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex //读写锁
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 1000),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *wsServer) Router(router *router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value string) {
	w.propertyLock.Lock() //加锁
	defer w.propertyLock.Unlock()
	w.property[key] = value
}
func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.RLock()
	defer w.propertyLock.RUnlock()
	return w.property[key], nil
}
func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}
func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}
func (w *wsServer) Push(name string, data string) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

// 通道一旦建立，那么收发消息就得要一直监听才行
func (w *wsServer) Start() {
	//启动读取数据的处理逻辑
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) readMsgLoop() {

	defer func() {
		if err := recover(); err != nil {
			log.Fatal(err)
			w.Close()
		}
	}()
	//先读到客户端发送过来的数据，然后进行处理，然后再回消息
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误:", err)
			break
		}
		//收到消息，解析消息。前端发送过来的消息就是json格式
		//1.data做一个解压，unzip
		//2.前端的消息，加密消息。进行解密
		//3.data转为body
		/*		data, err := utils.Decompress(data) //解压缩消息
				if err != nil {
					log.Println("解压缩数据出错，非法格式：", err)
					continue
				}

				body := &ReqBody{}*/
		fmt.Println(data)
	}
	w.Close()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			fmt.Println(msg)
		}
	}
}

func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}
