package controller

import (
	"newGoSang/net"
	"newGoSang/server/login/proto"
)

var DefaultAccount = &Account{}

type Account struct {
}

func (a *Account) Router(r *net.Router) {
	g := r.Group("account")
	g.AddRouter("login", a.login)
}

func (a *Account) login(req *net.WsMsgReq, rsp *net.WsMsgRsp) {
	//根据用户名，密码 查询user表
	//密码比对
	//保存用户的登录记录
	//保存用户最后一次登录信息
	//客户端需要一个session  jwt生成加密算法
	//客户端访问的时候，判断用户是否合法
	rsp.Body.Code = 0
	loginRes := &proto.LoginRsp{}
	loginRes.UId = 1
	loginRes.Username = "admin"
	loginRes.Session = "as"
	loginRes.Password = ""
	rsp.Body.Msg = loginRes
}
