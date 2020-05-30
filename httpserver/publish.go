package httpserver

import (
	"encoding/json"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"nsp/nspserver"
	"nsp/tools/log"
	"time"
)

func Send(ctx *gin.Context)  {
	topic := ctx.Query("topic")
	data :=  ctx.Query("data")
	log.Instance().Info("info",zap.Any("data",data))

	var dat interface{}
	if err := json.Unmarshal([]byte(data), &dat); err != nil {
		req := Response{Code: 201, Msg: err.Error()}
		ctx.JSON(200, req)
		return
	}
	req := Response{Code: 201,Ts: time.Now().Unix(),Msg: "发送成功"}
	ctx.JSON(200, req)

	msg:=nspserver.Response{Code: 200,Topic: topic,Ts: time.Now().Unix(),Data: dat}
	nspserver.Pub.SendIndex(topic,msg)

	return


}