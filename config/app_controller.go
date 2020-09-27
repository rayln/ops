package config

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/controller"
)

func InitController(application *mvc.Application, app *iris.Application, engine *xorm.EngineGroup) {
	//添加一个Controller用例，可以参考进行使用
	application.Handle(new(controller.TestController))
	//第二个，第三个。。。。
}
