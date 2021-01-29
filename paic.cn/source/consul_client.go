package source

import (
	"fmt"
	"log"
	"net"
	"strconv"

	"github.com/hashicorp/consul/api"
)

// GetService 根据服务名获取注册中心地址
func GetService(service string) map[string]struct{} {
	var lastIndex uint64

	config := api.DefaultConfig()
	// consul 注册中心地址
	config.Address = "127.0.0.1:8500"
	// consul client
	client, err := api.NewClient(config)

	if err != nil {
		fmt.Println("api new client is failed, err:", err)
		return nil
	}
	// 获取地址
	services, metainfo, err := client.Health().Service(service, "", true, &api.QueryOptions{
		WaitIndex: lastIndex, // 同步直到有新的更新
	})
	if err != nil {
		log.Println("error retrieving instances from Consul:", err)
	}
	lastIndex = metainfo.LastIndex
	addrs := map[string]struct{}{}
	// 返回服务地址
	for _, service := range services {
		log.Println("service.Service.Address:", service.Service.Address, "service.Service.Port:", service.Service.Port)
		addrs[net.JoinHostPort(service.Service.Address, strconv.Itoa(service.Service.Port))] = struct{}{}
	}
	return addrs
}
