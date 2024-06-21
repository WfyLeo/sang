package net

import "strings"

type HandlerFunc func(req *WsMsgReq, rsp *WsMsgRsp)

// account login//login
type group struct {
	prefix     string
	handlerMap map[string]HandlerFunc
}

func (g *group) AddRouter(name string, handlerFunc HandlerFunc) {
	g.handlerMap[name] = handlerFunc
}

func (r *Router) Group(prefix string) *group {
	g := &group{
		prefix:     prefix,
		handlerMap: make(map[string]HandlerFunc),
	}
	r.group = append(r.group, g)
	return g
}

func (g *group) exec(name string, req *WsMsgReq, rsp *WsMsgRsp) {
	handlerFunc := g.handlerMap[name]
	if handlerFunc != nil {
		handlerFunc(req, rsp)
	}
}

type Router struct {
	group []*group
}

func NewRouter() *Router {
	return &Router{}
}

func (r *Router) Run(req *WsMsgReq, rsp *WsMsgRsp) {
	//req.Body.Name 这个就是路径
	strs := strings.Split(req.Body.Name, ".")
	prefix := ""
	name := ""
	if len(strs) == 2 {
		prefix = strs[0]
		name = strs[1]
	}
	for _, g := range r.group {
		if g.prefix == prefix {
			g.exec(name, req, rsp)
		}
	}
}
