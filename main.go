package main

import (
	"go.uber.org/zap"

	"go-kunpeng/library/log"
	"go-kunpeng/server"
	"go-kunpeng/service"
)

func main() {
	// 注册logger
	defer log.Logger.Sync()
	zap.ReplaceGlobals(log.Logger)

	// 启动canal服务
	go service.StartCanalClient()

	// 启动grpc接口服务
	server.Start()
}
