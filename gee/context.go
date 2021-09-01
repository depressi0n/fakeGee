package gee

import (
	"encoding/json"
	"fmt"
	"net/http"
)

// 设计Context的必要性，避免繁杂重复代码，高效构造HTTP响应；同时提供了解析动态路由的能力，将复杂性和扩展性保留在内部，对外简化接口
// Context作为参数在实现路由处理函数以及中间件函数中调用

// H 别名
type H map[string]interface{}

// Context 提供快速构造Query和PostForm参数方法，并提供快速构造String/Data/JSON/HTML响应对方法
type Context struct {
	Writer http.ResponseWriter
	Req    *http.Request

	Path   string
	Method string
	Params map[string]string

	StatusCode int
	// for middlewares
	handlers []HandlerFunc
	index int

	engine *Engine
}

func newContext(w http.ResponseWriter, req *http.Request) *Context {
	return &Context{
		Writer: w,
		Req:    req,
		Path:   req.URL.Path,
		Method: req.Method,
		index:-1, // 记录当前执行的中间件
	}
}

func (c *Context) Next() {
	c.index++
	s:=len(c.handlers)
	for ; c.index<s;c.index++{
		c.handlers[c.index](c)
	}
}

func (c *Context) Fail(code int,err string) {
	c.index=len(c.handlers)
	c.JSON(code,H{"message":err})
}

func (c *Context) Param(key string) string {
	return c.Params[key]
}

func (c *Context) PostForm(key string) string {
	return c.Req.FormValue(key)
}

func (c *Context) Query(key string) string {
	return c.Req.URL.Query().Get(key)
}
func (c *Context) Status(code int) {
	c.StatusCode = code
	c.Writer.WriteHeader(code)
}
func (c *Context) SetHeader(key string, value string) {
	c.Writer.Header().Set(key, value)
}
func (c *Context) String(code int, format string, values ...interface{}) {
	c.SetHeader("Content-Type", "text/plain")
	c.Status(code)
	c.Writer.Write([]byte(fmt.Sprintf(format, values...)))
}
func (c *Context) JSON(code int, obj interface{}) {
	c.SetHeader("Content-Type", "application/json")
	c.Status(code)
	encoder := json.NewEncoder(c.Writer)
	if err := encoder.Encode(obj); err != nil {
		http.Error(c.Writer, err.Error(), 500)
	}
}
func (c *Context) Data(code int, data []byte) {
	c.Status(code)
	c.Writer.Write(data)
}
func (c *Context) HTML(code int, name string, data interface{}) {
	c.SetHeader("Content-Type", "text/html")
	c.Status(code)
	// 支持模板文件名选择模板进行渲染
	if err:=c.engine.htmlTemplates.ExecuteTemplate(c.Writer,name,data);err!=nil{
		c.Fail(500,err.Error())
	}
}

