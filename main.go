package main

import (
	"fmt"
	"gee"
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
	r := gee.Default()
	r.GET("/", func(ctx *gee.Context) {
		ctx.String(http.StatusOK,"Hello Depressi0n")
	})
	r.GET("/panic", func(ctx *gee.Context) {
		names:=[]string{"depressi0n"}
		ctx.String(http.StatusOK,names[100])
	})
	r.RUN(":9999")
}
