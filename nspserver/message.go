package nspserver



//数据接收
type Request struct {
	Type  string      `json:"type"`  //订阅类型 subscribe 订阅  unsubscribe 取消订阅
	Topic   string     `json:"topic"`   //订阅通道
	Data  interface{} `json:"data"` //接收数据
}

//数据返回
type SubResponse struct {
	Code 		int 	`json:"code"`       //状态
	Tytp 		string `json:"tytp"`
	Topic    	string `json:"topic"`  //主题
	Ts     		int64  `json:"ts"`    //时间
}

//错误的提示
type ErrResponse struct {
	Code   int `json:"status"`
	Err    string `json:"err"`
	Ts     int64  `json:"ts"`
}

//推送响应消息
type Response struct {
	Code     	int 		`json:"code"`
	Topic   	string      `json:"topic"`
	Ts   		int64       `json:"ts"`
	Data 		interface{} `json:"data"`
}


