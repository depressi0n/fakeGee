package gee

import (
	"net/http"
	"strings"
)

// 前缀树根据HTTP请求路径用/分隔成多段作为一个节点，继承到router，提供动态路由功能之参数匹配，通配符功能
type router struct {
	roots    map[string]*node
	handlers map[string]HandlerFunc
}

func newRouter() *router {
	return &router{
		roots:    make(map[string]*node),
		handlers: make(map[string]HandlerFunc),
	}
}
func parsePatter(pattern string) []string {
	ss := strings.Split(pattern, "/")
	parts := make([]string, 0)
	for _, item := range ss {
		if item != "" {
			parts = append(parts, item)
			if item[0] == '*' {
				break
			}
		}
	}
	return parts
}
func (r *router) addRoute(method string, pattern string, handler HandlerFunc) {
	parts := parsePatter(pattern)
	key := method + "_" + pattern
	_, ok := r.roots[method]
	if !ok {
		r.roots[method] = &node{}
	}
	r.roots[method].insert(pattern, parts, 0)
	r.handlers[key] = handler
}
func (r *router) getRoute(method string, path string) (*node, map[string]string) {
	searchParts := parsePatter(path)
	params := make(map[string]string)
	root, ok := r.roots[method]
	if !ok{
		return nil,nil
	}
	n:=root.search(searchParts,0)
	if n!=nil{
		parts:=parsePatter(n.pattern)
		for index, part := range parts {
			if part[0]==':'{
				params[part[1:]]=searchParts[index]
			}
			if part[0]=='*' && len(part)>1{
				params[part[1:]]=strings.Join(searchParts[index:],"/")
				break
			}
		}
		return n,params
	}
	return nil,nil
}
func (r *router) handle(c *Context) {
	n, params := r.getRoute(c.Method, c.Path)
	if n!=nil{
		c.Params=params
		key:=c.Method+"_"+n.pattern
		// 将路由匹配到的Handler添加到c.handlers中，执行c.Next()
		c.handlers=append(c.handlers,r.handlers[key])
	}else{
		c.handlers=append(c.handlers,func(c *Context){
			c.String(http.StatusNotFound, "404 NOT FOUND: %s \n", c.Path)
		})
	}
	// 中间件编写过程中，可以增加前后处理，这样就可以在请求处理之前和之后添加处理阶段
	// 开始执行
	c.Next()
}
