package net

import (
	"encoding/json"
	"errors"
	"github.com/forgoer/openssl"
	"github.com/gorilla/websocket"
	"log"
	"newGoSang/utils"
	"sync"
)

// webSocket服务
type wsServer struct {
	wsConn       *websocket.Conn
	router       *Router
	outChan      chan *WsMsgRsp
	Seq          int64
	property     map[string]interface{}
	propertyLock sync.RWMutex
}

func NewWsServer(wsConn *websocket.Conn) *wsServer {
	return &wsServer{
		wsConn:   wsConn,
		outChan:  make(chan *WsMsgRsp, 100),
		property: make(map[string]interface{}),
		Seq:      0,
	}
}

func (w *wsServer) Router(router *Router) {
	w.router = router
}

func (w *wsServer) SetProperty(key string, value interface{}) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	w.property[key] = value
}

func (w *wsServer) GetProperty(key string) (interface{}, error) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	if value, ok := w.property[key]; ok {
		return value, nil
	} else {
		return nil, errors.New("no property found")
	}
}

func (w *wsServer) RemoveProperty(key string) {
	w.propertyLock.Lock()
	defer w.propertyLock.Unlock()
	delete(w.property, key)
}

func (w *wsServer) Addr() string {
	return w.wsConn.RemoteAddr().String()
}

func (w *wsServer) Push(name string, data interface{}) {
	rsp := &WsMsgRsp{Body: &RspBody{Name: name, Msg: data, Seq: 0}}
	w.outChan <- rsp
}

func (w *wsServer) Start() {
	go w.readMsgLoop()
	go w.writeMsgLoop()
}

func (w *wsServer) writeMsgLoop() {
	for {
		select {
		case msg := <-w.outChan:
			w.write(msg)
		}
	}
}

func (w *wsServer) write(msg *WsMsgRsp) {
	data, err := json.Marshal(msg.Body)
	if err != nil {
		log.Println(err)
	}
	secretKey, err := w.GetProperty("secretKey")
	if err == nil {
		//有加密
		key := secretKey.(string)
		//数据做加密
		data, _ = utils.AesCBCEncrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
	}
	if data, err := utils.Zip(data); err == nil {
		w.wsConn.WriteMessage(websocket.BinaryMessage, data)
	}
}

func (w *wsServer) readMsgLoop() {
	//先读到客户端发送过来的数据，然后进行处理，然后再回消息
	//经过路由，实习处理情绪
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
			w.Close()
		}
	}()
	for {
		_, data, err := w.wsConn.ReadMessage()
		if err != nil {
			log.Println("收消息出现错误")
			break
		}
		//收到消息，解析消息，前端发送过来的就是json格式
		//1.消息进行一个解压
		data, err = utils.UnZip(data)
		if err != nil {
			log.Println("解压数据出错，非法格式")
		}
		//2.消息进行解密
		secretKey, err := w.GetProperty("secretKey")
		if err == nil {
			//有加密
			key := secretKey.(string)
			//客户端传递过来的数据是加密的，需要解密
			d, err := utils.AesCBCDecrypt(data, []byte(key), []byte(key), openssl.ZEROS_PADDING)
			if err != nil {
				log.Println("数据格式有误，解密失败：", err)
				//出错后，捂手
				//w.Handshake()
			} else {
				data = d
			}

		}
		//3.data 转为 body
		body := &ReqBody{}
		err = json.Unmarshal(data, body)
		if err != nil {
			log.Println("数据格式有误，非法格式:", err)
		} else {
			//拿到了前端传递的数据，拿上这些数据，去具体的业务进行处理
			req := &WsMsgReq{Conn: w, Body: body}
			rsp := &WsMsgRsp{Body: &RspBody{Name: body.Name, Seq: req.Body.Seq}}
			w.router.Run(req, rsp)
			w.outChan <- rsp
		}
	}
	w.Close()
}

func (w *wsServer) Close() {
	_ = w.wsConn.Close()
}

const HandshakeMsg = "handshake"

// 当游戏客户端 发送请求的时候，会先进行握手协议
// 后端会发送对应的加密key给客户端
// 客户端再发送数据的时候，就会使用key进行一个加密处理
func (w *wsServer) Handshake() {
	secretKey := ""
	key, err := w.GetProperty("secretKey")
	if err == nil {
		secretKey = key.(string)
	} else {
		secretKey = utils.RandSeq(16)
	}

	handshake := &Handshake{Key: secretKey}

	body := &RspBody{Name: HandshakeMsg, Msg: handshake}

	if data, err := json.Marshal(body); err == nil {
		if secretKey != "" {
			w.SetProperty("secretKey", secretKey)
		} else {
			w.RemoveProperty("secretKey")
		}
		if data, err := utils.Zip(data); err == nil {
			w.wsConn.WriteMessage(websocket.BinaryMessage, data)
		}
	}
}
