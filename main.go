package main

import (
	"go-kunpeng/server"
	"go-kunpeng/service"
)

func main() {
	// 启动canal服务
	service.StartCanalClient()
	// 启动grpc接口服务
	server.Start()
}
