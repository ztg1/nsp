// Copyright 2019 GitBitEx.com
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package httpserver

import (
	"errors"
	"github.com/gin-gonic/gin"
	"net/http"
	"nsp/config"
)

func CheckToken() gin.HandlerFunc {
	return func(c *gin.Context) {
		//ip处理
		clientIp:=c.ClientIP()
		if(IsIp(clientIp)==false){
			c.AbortWithStatusJSON(http.StatusForbidden, newMessageVo(201,errors.New("Not ip")))
			return
		}
		c.Next()
	}

}

//判断ip是否合法
func IsIp(clientIp string) bool  {
	ips:=config.GetConfig().Ips
	if(len(ips)>0){
		for _,index:=range ips{
			if(index == clientIp){

				return  true
			}
		}
		return  false
	}
	return true
}