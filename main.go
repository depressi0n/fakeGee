package main

import (
	"fmt"
	"gee"
	"html/template"
	"net/http"
	"time"
)

type student struct {
	Name string
	Age int8
}
func FormatAsDate(t time.Time)string{
	year,month,day:=t.Date()
	return fmt.Sprintf("%d-%02d-%02d",year,month,day)
}

func main() {
	r := gee.New()
	r.Use(gee.Logger()) // 全局中间件，最顶层
	r.SetFuncMap(template.FuncMap{
		"FormatAsDate":FormatAsDate,
	})
	r.LoadHTMLGlob("templates/*")
	r.Static("/assets","./static") // 静态文件绑定

	stu1:=&student{
		Name: "Depressi0n",
		Age:  22,
	}
	stu2:=&student{
		Name: "Ocean",
		Age:  21,
	}

	r.GET("/", func(ctx *gee.Context) {
		ctx.HTML(http.StatusOK,"css.tmpl",nil)
	})
	r.GET("/student",func(c *gee.Context){
		c.HTML(http.StatusOK,"arr.tmpl",gee.H{
			"title":"gee",
			"stuArr":[2]*student{stu1,stu2},
		})
	})
	r.GET("/date",func(c *gee.Context){
		c.HTML(http.StatusOK,"arr.tmpl",gee.H{
			"title":"gee",
			"now":time.Date(2021,9,1,0,0,0,0,time.UTC),
		})
	})
	r.RUN(":9999")
}
