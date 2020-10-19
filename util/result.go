package util

import (
	"encoding/json"
	"github.com/kataras/iris/v12/mvc"
)

type Result struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (that *Result) Response() mvc.Result {
	if that.Code != 0 {
		that.Message = "网络出现异常，请重新再试！"
	}
	return mvc.Response{
		//Object: StructToMap(*that),
		Object: that,
	}
}

func (that *Result) ResponseWs() string {
	if that.Code != 0 {
		that.Message = "网络出现异常，请重新再试！"
	}
	temp, err := json.Marshal(that)
	if err != nil {
		panic(err)
	}
	return string(temp)
}
