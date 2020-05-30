package init

import (
	"nsp/config"
	"nsp/tools/log"
)

func init()  {
	log.Init("snp")
	config.GetConfig()

}
