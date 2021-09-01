package main

import (
	"gee"
	"log"
	"net/http"
	"time"
)

func onlyForV2() gee.HandlerFunc{
	return func(ctx *gee.Context) {
		t:=time.Now()
		ctx.Fail(500,"Internal Server Error")
		log.Printf("[%d] %s in %v for group v2",ctx.StatusCode,ctx.Req.RequestURI,time.Since(t))
	}
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // 全局中间件，最顶层
	r.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK,"<h1>Hello Gee</h1>")
	})
	r.GET("/index", func(c *gee.Context) {
		c.HTML(http.StatusOK, "<h1>Index Page</h1>")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/", func(c *gee.Context) {
			c.HTML(http.StatusOK, "<h1>Hello Gee</h1>")
		})
		v1.GET("/hello", func(c *gee.Context) {
			// expect /hello/?name=depressi0n
			c.String(http.StatusOK, "hello %s,you're at %s\n", c.Query("name"), c.Path)
		})
	}
	v2 := r.Group("/v2")
	v2.Use(onlyForV2())
	{
		v2.GET("/hello/:name", func(c *gee.Context) {
			// expect /hello/depressi0n
			c.String(http.StatusOK, "hello %s,you're at %s\n", c.Param("name"), c.Path)
		})
		v2.GET("/assets/*filepath", func(c *gee.Context) {
			c.JSON(http.StatusOK, gee.H{"filepath": c.Param("filepath")})
		})
	}

	r.RUN(":9999")
}
