package init

import (
	"nsp/httpserver"
	"nsp/nspserver"
)

func BootStart()  {

	nspserver.StartServer()//websocet 服务启动
    httpserver.StartHttpServer()   //http启动

}