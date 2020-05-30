package nspserver

import "sync"

type subscribe struct {
	subscribers map[string]map[string]*Client //客户端订阅
	mu          sync.RWMutex                  //互斥锁
}

//实例化 订阅结构体
func newSubscription() *subscribe {
	return &subscribe{subscribers: map[string]map[string]*Client{}}
}

//订阅那个市场
func (s *subscribe) subscribe(channel string, client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.subscribers[channel]
	if !found {
		s.subscribers[channel] = map[string]*Client{}
	}

	_, found = s.subscribers[channel][client.id]
	if found {
		return false
	}
	s.subscribers[channel][client.id] = client
	return true
}

//取消那个市场订阅
func (s *subscribe) unsubscribe(channel string, client *Client) bool {
	s.mu.Lock()
	defer s.mu.Unlock()

	_, found := s.subscribers[channel]
	if !found {
		return false
	}

	_, found = s.subscribers[channel][client.id]
	if !found {
		return false
	}
	delete(s.subscribers[channel], client.id)
	return true
}

//后台推送消息给订阅该topic的client
//topic 主题
//msg 内容
func (s *subscribe) Publish(topic string, msg interface{}) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	_, found := s.subscribers[topic]
	if !found {

		return
	}
	for _, c := range s.subscribers[topic] {
		c.writeCh <- msg
	}
}