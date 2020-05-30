package config

import (
	"encoding/json"
	"go.uber.org/zap"
	"io/ioutil"
	"nsp/tools/log"

	"sync"
)

var config GbConfig
var configOnce   sync.Once

func GetConfig() *GbConfig  {
	configOnce.Do(func() {
		bytes, err := ioutil.ReadFile("./etc/conf.json")
		if err != nil {
			log.Instance().Panic("panic",zap.Any("conf.json-err",err))
		}

		err = json.Unmarshal(bytes, &config)
		if err != nil {
			log.Instance().Panic("panic",zap.Any("json-unmarshal",err))
		}
	})
	return &config
}


//配置文件
type GbConfig struct {
   HttpServerc   HttpServerConfig  `json:"http_serverc"`
   WebSocketc    WebSocketConfig  	`json:"web_socketc"`
   Ips           []string  			`json:"ips"`
   Debug         bool				`json:"debug"`
}

//http服务访问地址
type HttpServerConfig struct {
	Addr   string `json:"addr"`
}


//websocket 访问地址
type WebSocketConfig struct {
	Addr	string	`json:"addr"`
	Path 	string 	`json:"path"`
}

