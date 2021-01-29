////+build wireinject
package main

import (
	"fmt"
	"log"
	"net/http"

	consulapi "github.com/hashicorp/consul/api"
	"paic.cn/source"
)

func main() {
	fmt.Println("hello world!")
	config := consulapi.DefaultConfig()
	fmt.Println(config)
	sers := source.GetService("java-demo")
	fmt.Println("services:", sers)
	http.HandleFunc("/", source.RouteHandler)
	log.Fatal(http.ListenAndServe(":2002", nil))
}
