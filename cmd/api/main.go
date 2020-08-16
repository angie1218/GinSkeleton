package main

import (
	"goskeleton/app/utils/config"
	_ "goskeleton/bootstrap"
	"goskeleton/routers"
)

// 这里可以存放门户类网站入口
func main() {
	routers := routers.InitApiRouter()
	routers.Run(config.CreateYamlFactory().GetString("HttpServer.Api.Port"))
}