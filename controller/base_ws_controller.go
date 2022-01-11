package controller

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/entity"
	"runtime"
	"xorm.io/xorm"
)

type BaseWsController struct {
	entity.BaseWsEntity
	Request context.Context
}

/**
开启事务
*/
func (that *BaseWsController) Begin() *entity.BaseEntity {
	var entitys = new(entity.BaseEntity)
	if entitys.Save == nil {
		//创建事务
		entitys.Save = that.Engine.NewSession()
		//开启事务
		entitys.Save.Begin()
		entitys.Logger, entitys.App, entitys.Engine, entitys.Redis = that.Logger, that.App, that.Engine, that.Redis
	}
	if entitys.Load == nil {
		entitys.Load = that.Engine.Slave()
	}
	if len(that.EngineOther) > 0 {
		entitys.EngineOther = that.EngineOther
		entitys.SaveOther = make([]*xorm.EngineGroup, len(that.EngineOther))
		entitys.LoadOther = make([]*xorm.Engine, len(that.EngineOther))
	}
	for i := 0; i < len(that.EngineOther); i++ {
		entitys.SaveOther[i] = that.EngineOther[i]
		entitys.LoadOther[i] = that.EngineOther[i].Slave()
	}
	return entitys
}

/**
提交事务
*/
func (that *BaseWsController) Commit(entitys *entity.BaseEntity) {
	if entitys.Save != nil {
		//提交事务
		entitys.Save.Commit()
	}
}
func (that *BaseWsController) Close(entitys *entity.BaseEntity) {
	if entitys.Save != nil {
		//关闭事务
		entitys.Save.Close()
		entitys.Save = nil
	}
}

/**
开始执行业务逻辑
事务开启和异常捕获
*/
/*func (that *BaseController) Start(request context.Context, baseEntity *entity.BaseEntity, serviceFunc func(*entity.BaseEntity)mvc.Result) mvc.Result{
	defer that.handleException(request)
	that.Begin()
	result := serviceFunc(baseEntity)
	that.Commit()
	return result
}*/
func (that *BaseWsController) Start(serviceFunc func(*entity.BaseEntity) string) (result string) {
	//TODO update
	result = "{\"code\":1,\"message\":\"system error!\",\"data\":\"\",\"system_error\":1}"
	var enti = that.Begin()
	defer that.Close(enti)
	defer that.handleException(&result, enti)
	//传入entity到用户中。然后再做新的
	result = serviceFunc(enti)
	that.Commit(enti)
	return result
}

func (that *BaseWsController) BeforeActivation(a mvc.BeforeActivation) {
	//拦截链接，主要用于Session，开启一个Session，并且捕获异常进行回滚
	a.Router().Use(func(context context.Context) {
		context.Next()
	})
	a.Router().Done(func(context context.Context) {
		//跳转结束
		context.Next()
	})
}
func (that *BaseWsController) AfterActivation(a mvc.AfterActivation) {

}

/**
* @Description: 错误信息处理
 */
func (that *BaseWsController) exceptionRecover(err interface{}) *string {
	var stacktrace string
	for i := 1; ; i++ {
		_, f, l, got := runtime.Caller(i)
		if !got {
			break
		}
		stacktrace += fmt.Sprintf("%s:%d\n", f, l)
	}

	errMsg := fmt.Sprintf("错误信息: %s", err)
	// when stack finishes
	logMessage := fmt.Sprintf("从错误中回复：\n")
	logMessage += errMsg + "\n"
	logMessage += fmt.Sprintf("\n%s", stacktrace)
	that.Logger.Error(logMessage)
	// 打印错误日志
	// 返回错误信息
	temp := "{data: \"\", code: 1, message: \"The server is error. Please try again at a moment!\",\"system_error\":1}"
	return &temp
}

/**
处理异常信息
*/
func (that *BaseWsController) handleException(result *string, entitys *entity.BaseEntity) {
	if err := recover(); err != nil {
		if entitys.Save != nil {
			//事务回滚
			entitys.Save.Rollback()
		}
		//异常处理
		result = that.exceptionRecover(err)
	} else {
		temp := "{data: \"handleException\", code: 0, message: \"\",,\"system_error\":0}"
		result = &temp
	}
}
