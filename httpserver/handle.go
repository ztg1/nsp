package httpserver

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"nsp/tools/log"
	"time"
)

type HttpServer struct {
    
	addr string
}

func NewHttpServer(addr string)*HttpServer  {

	return &HttpServer{addr: addr}
}

func (server *HttpServer) Start() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard

	r := gin.Default()
	r.Use(setCROSOptions)
	r.GET("/", func(ctx *gin.Context) {
		times := time.Now().Unix()
		msg := Response{Code: 200, Ts: times,Data: "ok"}
		ctx.JSON(http.StatusOK, msg)
		return
	})
	r.GET("/api/getname", func(ctx *gin.Context) {
		times := time.Now().Unix()
		msg := Response{Code: 200, Ts: times}
		ctx.JSON(http.StatusOK, msg)
		return
	})



	private := r.Group("/", CheckToken())
	{
		private.POST("/api/send", Send)                //测试接口
	}
	err := r.Run(server.addr)
	if err != nil {
		log.Instance().Panic("panic", zap.Any("httpServer-err", err))
		return
	}
}

func setCROSOptions(c *gin.Context) {
	method := c.Request.Method
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Access-Control-Allow-Headers", "*")
	c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
	c.Header("Content-Type", "application/json")
	//放行所有OPTIONS方法
	if method == "OPTIONS" {
		c.AbortWithStatus(http.StatusNoContent)
	}

}

