package httpserver

import "time"

//响应数据....
type Requet struct {

}

//说返回
func newMessageVo(code int, error error) *Response {
	return &Response{
		Code:    code,
		Msg:     error.Error(),
		Ts:		time.Now().Unix(),
	}
}

//响应的
type Response struct {
	Code   	int 		`json:"code"`
	Msg     string    	`json:"msg"`
	Ts      int64		`json:"ts"`
	Data  	interface{}  `json:"data"`
}
