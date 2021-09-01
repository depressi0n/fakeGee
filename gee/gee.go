package gee

import (
	"net/http"
)

// 使用一个自定义的结构体取代标准库的默认实例，相当于一个拦截器，将所有的HTTP请求全部交与这个结构体处理
// 这样带来的好处是，统一了控制入口，并且可以自定义路由规则，统一处理逻辑如日志、异常处理等

type Engine struct {
	router *router
}
type HandlerFunc func(ctx *Context)

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	c:=newContext(w,req)
	engine.router.handle(c)
}
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	engine.router.addRoute(method,pattern,handler)
}
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}

//RUN 运行Server
func (engine *Engine) RUN(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// New 创建gee实例，使用GET方法添加路由，使用RUN启动Web服务
func New() *Engine {
	return &Engine{router: newRouter()}
}

