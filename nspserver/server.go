package nspserver

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
	"io/ioutil"
	"net/http"
	"nsp/tools/log"
)

//服务结构体
type Server struct {
	addr string //服务器地址
	path string //后面地址
	sub  *subscribe
}

//实例化 地址订阅消息
func NewServer(addr, path string, sub *subscribe) *Server {
	return &Server{
		addr: addr,
		path: path,
		sub:  sub,
	}
}

//示例话clien端
func (s *Server) ws(c *gin.Context) {

	upGrader := websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
	conn, err := upGrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Instance().Debug("err", zap.Any("ws-err", err))
		return
	}
	NewClient(conn, s.sub).startServe() //-------------
}

//路由 /ws 进来的
func (s *Server) Run() {
	gin.SetMode(gin.ReleaseMode)
	gin.DefaultWriter = ioutil.Discard
	r := gin.Default()
	r.GET(s.path, s.ws)
	err := r.Run(s.addr)

	if err != nil {
		log.Instance().Panic("panic", zap.Any("panic-ws", err))
		return
	}

}
