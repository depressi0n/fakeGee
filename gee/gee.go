package gee

import (
	"fmt"
	"net/http"
)

// 使用一个自定义的结构体取代标准库的默认实例，相当于一个拦截器，将所有的HTTP请求全部交与这个结构体处理
// 这样带来的好处是，统一了控制入口，并且可以自定义路由规则，统一处理逻辑如日志、异常处理等

type Engine struct {
	router map[string]HandlerFunc
}
type HandlerFunc func(http.ResponseWriter, *http.Request)

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	key:=req.Method+"_"+req.URL.Path
	if handler,ok:=engine.router[key];ok{
		handler(w,req)
	}else{
		fmt.Fprintf(w, "404 NOT FOUND: %s \n", req.URL)
	}
}
func (engine *Engine) addRoute(method string, pattern string, handler HandlerFunc) {
	key := method + "_" + pattern
	engine.router[key] = handler
}
func (engine *Engine) GET(pattern string, handler HandlerFunc) {
	engine.addRoute("GET", pattern, handler)
}

func (engine *Engine) POST(pattern string, handler HandlerFunc) {
	engine.addRoute("POST", pattern, handler)
}
func (engine *Engine) RUN(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// New 创建gee实例，使用GET方法添加路由，使用RUN启动Web服务
func New() *Engine {
	return &Engine{router: make(map[string]HandlerFunc)}
}
