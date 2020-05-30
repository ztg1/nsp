package nspserver

import (
	"go.uber.org/zap"
	"nsp/config"
	"nsp/tools/log"
)

func StartServer()  {

	gbeConfig := config.GetConfig()
	sub := newSubscription()
	NewPublish(sub) //实例化推送
	go NewServer(gbeConfig.WebSocketc.Addr, gbeConfig.WebSocketc.Path, sub).Run() //实例化主题
	log.Instance().Info("info",zap.Any("info","websocket-server-start:"+gbeConfig.WebSocketc.Addr+gbeConfig.WebSocketc.Path))
}
