package main

import (
	"go.uber.org/zap"
	."nsp/init"
	"nsp/tools/log"
)


func main()  {
	BootStart()  //启动服务
	log.Instance().Info("info", zap.Any("srv-start", "server start ........."))
	select {}
}
