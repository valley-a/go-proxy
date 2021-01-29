package main

import (
	"log"
	"net/http"

	"paic.cn/source"
)

func main() {
	log.Println("agent start!")
	// 配置路由 http handler
	http.HandleFunc("/", source.RouteHandler)
	// 启动监听并设置端口
	log.Fatal(http.ListenAndServe(":2002", nil))
}
