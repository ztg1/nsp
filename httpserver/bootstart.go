package httpserver

import (
	"go.uber.org/zap"
	"nsp/config"
	"nsp/tools/log"
)

func StartHttpServer() {
	httpServer := NewHttpServer(config.GetConfig().HttpServerc.Addr)
	go httpServer.Start()
	log.Instance().Info("info", zap.Any("httpServer", "server start..启动成功 端口:"+config.GetConfig().HttpServerc.Addr))

}
