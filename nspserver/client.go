package nspserver

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	uuid "github.com/satori/go.uuid"
	"go.uber.org/zap"
	"nsp/tools/log"
	"sync"
	"time"
)

const (
	writeWait      = 10 * time.Second
	pongWait       = 60 * time.Second
	pingPeriod     = (pongWait * 9) / 10
	maxMessageSize = 512 * 2
)

type Client struct {
	id   		string
	conn 		*websocket.Conn
	writeCh  	chan interface{} //用于写数据的
	sub         *subscribe
	channels 	map[string]struct{}
	pingTime 	int64 //由于记录ping值的时间
	mu       	sync.Mutex
}

//实例化
func NewClient(conn *websocket.Conn, sub *subscribe) *Client {
	u1:= uuid.NewV4()
	return &Client{
		id:       u1.String(),
		conn:     conn,
		writeCh:  make(chan interface{}, 256*2), //提示消息返回个客户端
		sub:      sub,
		pingTime: time.Now().Unix(), //
		channels: map[string]struct{}{},
	}
}

func (c *Client) startServe() {
	go c.runReader()
	go c.runWriter()
}

//写入数据
func (c *Client)runWriter()  {

	ticker:=time.NewTicker(pingPeriod) //启动一个定时器
	defer func() {
		ticker.Stop()
		_=c.conn.Close()
	}()


	for {
		select {
		case message := <-c.writeCh:
			err := c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err != nil {
				_ = c.conn.WriteMessage(websocket.CloseMessage, []byte{})
				c.close()
				return
			}

			buf, err := json.Marshal(message)
			if err != nil {
				continue
			}
			err = c.conn.WriteMessage(websocket.TextMessage, buf)
			if err != nil {
				c.close()
				return
			}

		case <-ticker.C: //监听客户是否在线 不在线则把客户端踢掉
			_ = c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			err := c.conn.WriteMessage(websocket.PingMessage, nil)
			if err != nil {
				c.close()
				return
			}
			ping := struct {
				Type string `json:"type"`
			}{Type: "ping"}
			c.writeCh <- ping
			//如果54秒内没有发送pong  就把客户端踢掉
			if time.Now().Unix()-c.pingTime > 60 {
				c.close()
				return
			}

		}
	}

}

//读取数据
func (c *Client) runReader() {
	c.conn.SetReadLimit(maxMessageSize)
	err := c.conn.SetReadDeadline(time.Now().Add(pongWait))
	if err != nil {
		log.Instance().Error("err", zap.Any("client-red-err", err))
		
	}
	c.conn.SetPongHandler(func(string) error {

		return c.conn.SetReadDeadline(time.Now().Add(pongWait))
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			c.close()
			log.Instance().Debug("err", zap.Any("debug-red-err", err))
			break
		}
		var req Request
		err = json.Unmarshal(message, &req)
		if err != nil {
			log.Instance().Debug("err", zap.Any("client-red-err", string(message)))
			c.close()
			break
		}

		c.onMessage(&req)
	}
}

//消息处理
func (c *Client)onMessage(req *Request)  {

	switch req.Type {
	case "subscribe": //主题订阅
		c.onSub(req.Topic)
	case "unsubscribe": //取消主题订阅
		c.UnSub(req.Topic)
	case "pong":
		c.pingTime=time.Now().Unix()

	default:
		msg := ErrResponse{Code: 300, Err: "type is subscribe or unsubscribe or pong", Ts: time.Now().Unix()}
		c.writeCh <- msg
	}
}

/**
  加入订阅组

*/

func (c *Client) onSub(topic string) {
	ok := c.subscribe(topic)
	if ok == true {
		req := SubResponse{Code: 200, Topic: topic, Ts: time.Now().Unix(),Tytp: "subscribe"}
		c.writeCh <- req
	}
}

//取消订阅组
func (c *Client) UnSub(topic string) {
	c.unsubscribe(topic)
	req := SubResponse{Code: 200, Topic: topic, Ts: time.Now().Unix(),Tytp: "unsubscribe"}
	c.writeCh <- req
}


//订阅消息
func (c *Client) subscribe(channel string) bool {
	c.mu.Lock()
	defer c.mu.Unlock()

	_, found := c.channels[channel]

	if found {
		return false
	}
	if c.sub.subscribe(channel, c) {
		c.channels[channel] = struct{}{}
		return true
	}
	return false
}

//取消订阅
func (c *Client) unsubscribe(channel string) {
	c.mu.Lock()
	defer c.mu.Unlock()

	if c.sub.unsubscribe(channel, c) {
		delete(c.channels, channel)
	}
}


//关闭客户端
func (c *Client) close() {
	c.mu.Lock()
	defer c.mu.Unlock()

	for channel := range c.channels {
		c.sub.unsubscribe(channel, c)
	}
}