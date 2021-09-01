package gee

import (
	"log"
	"net/http"
	"strings"
)

// 使用一个自定义的结构体取代标准库的默认实例，相当于一个拦截器，将所有的HTTP请求全部交与这个结构体处理
// 这样带来的好处是，统一了控制入口，并且可以自定义路由规则，统一处理逻辑如日志、异常处理等

type Engine struct {
	*RouterGroup
	router *router
	groups []*RouterGroup // 存储所有的分组，将所有和路由的函数全部交给RouterGroup实现
}
type HandlerFunc func(ctx *Context)

func (engine Engine) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// 首先判断请求用于哪些中间件，这里只使用前缀判断
	var middlewares []HandlerFunc
	for _,group:=range engine.groups{
		if strings.HasPrefix(req.URL.Path,group.prefix){
			middlewares=append(middlewares,group.middlewares...)
		}
	}
	c := newContext(w, req)
	c.handlers=middlewares
	engine.router.handle(c)
}

//RUN 运行Server
func (engine *Engine) RUN(addr string) error {
	return http.ListenAndServe(addr, engine)
}

// New 创建gee实例，使用GET方法添加路由，使用RUN启动Web服务
func New() *Engine {
	engine := &Engine{router: newRouter()}
	engine.RouterGroup = &RouterGroup{engine: engine}
	engine.groups = []*RouterGroup{engine.RouterGroup} // 顶层的GroupRouter
	return engine
}

type RouterGroup struct {
	prefix      string
	middlewares []HandlerFunc // 支持中间件
	parent      *RouterGroup  // 为了支持嵌套
	engine      *Engine       // 所有的组共享一个Engine实例
}

func (group *RouterGroup) Group(prefix string) *RouterGroup {
	engine := group.engine
	newGroup := &RouterGroup{
		prefix: group.prefix + prefix,
		parent: group,
		engine: engine,
	}
	engine.groups=append(engine.groups,newGroup)
	return newGroup
}

// addRoute 将Engine继承了RouterGroup的所有属性和方法，可以通过engine.addRoute添加路由，也可以通过分组添加路由
func (group *RouterGroup) addRoute(method string, comp string, handler HandlerFunc) {
	pattern:=group.prefix+comp
	log.Printf("Route %4s - %s",method,pattern)
	group.engine.router.addRoute(method, pattern, handler)
}

func (group *RouterGroup) GET(pattern string, handler HandlerFunc) {
	group.addRoute("GET", pattern, handler)
}

func (group *RouterGroup) POST(pattern string, handler HandlerFunc) {
	group.addRoute("POST", pattern, handler)
}
// Use 将中间件应用到Group
func (group *RouterGroup) Use(middlewares ...HandlerFunc) {
	group.middlewares=append(group.middlewares,middlewares...)
}