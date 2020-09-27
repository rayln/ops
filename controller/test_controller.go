package controller

import (
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/entity"
	"github.com/rayln/ops/service"
)

type TestController struct {
	BaseController
	Request     iris.Context
	TestService service.TestService
}

/**
页面参数，通过post的raw传过来json
*/
type Test struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

/***************************逻辑处理开始*******************************/
/**
POST请求。http://localhost:8888/test 请求。raw获取json
*/
func (that *TestController) PostTest() mvc.Result {
	var farmTestInfo entity.FarmTestInfo
	that.Request.ReadJSON(&farmTestInfo)
	//开启事务代理，固定写法
	return that.Start(that.Request, func() mvc.Result {
		return that.TestService.Test(&that.BaseEntity)
	})
}
