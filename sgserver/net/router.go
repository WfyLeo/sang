package net

type HandlerFunc func()

// 路由组
type group struct {
	prefix     string //组的前缀
	handlerMap map[string]HandlerFunc
}

// 路由
type router struct {
	group []*group
}
