package net

// 请求的格式
type ReqBody struct {
	Seq   int64  `json:"seq"`
	Name  string `json:"name"`
	Msg   string `json:"msg"`
	Proxy string `json:"proxy"`
}

// 返回的格式
type RspBody struct {
	Seq  int64  `json:"seq"`
	Name string `json:"name"`
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

type WsMsgReq struct {
	Body *ReqBody
	Conn WSConn
}

type WsMsgRsp struct {
	Body *RspBody
}

// 处理属性，比如request请求会有参数，请求中取参数，放参数
type WSConn interface {
	SetProperty(key string, value string)
	GetProperty(key string) (interface{}, error)
	RemoveProperty(key string)
	Addr() string
	Push(name string, data string)
}
