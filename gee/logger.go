package gee

import (
	"log"
	"time"
)

// middlewares 非业务技术类组件，提供一个插口，允许用户自定义功能，嵌入框架中
// 首先插入点不能太底层，否则中间件逻辑会很复杂，也不能太顶层，那样用户可以直接定义一组函数，失去了框架的优势
// 中间件的输入决定了扩展能力，暴露的参数太少则用户发挥空间有限

// 这里的中间件处理的输入是Context，插入点是框架接收到请求初始化Context后，允许用户使用自己定义的中间件做一些额外的处理，如记录日志等，以及对Context进行二次加工
// 允许等待Handler处理之后做一些额外操作如计算处理时间等

func Logger() HandlerFunc {
	return func(c *Context){
		t:=time.Now()
		c.Next()
		log.Printf("[%d] %s in %v",c.StatusCode,c.Req.RequestURI,time.Since(t))
	}
}
