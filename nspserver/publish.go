package nspserver

import "sync"


var Pub *Publish

type Publish struct {
	sub   *subscribe
	Mutex sync.Mutex
}

func NewPublish(sub *subscribe)  {
  Pub=&Publish{sub: sub,Mutex: sync.Mutex{}}
}

//发送消息给主题
func (p * Publish)SendIndex(topic string,msg interface{})  {
	p.sub.Publish(topic,msg)
}