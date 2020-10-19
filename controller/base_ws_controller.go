package controller

import (
	"fmt"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/entity"
	"runtime"
)

type BaseWsController struct {
	entity.BaseEntity
	Request context.Context
}

/**
开启事务
*/
func (that *BaseWsController) Begin() {
	if that.Save == nil {
		//创建事务
		that.Save = that.Engine.NewSession()
		//开启事务
		that.Save.Begin()
	}
	if that.Load == nil {
		that.Load = that.Engine.Slave()
	}
}

/**
提交事务
*/
func (that *BaseWsController) Commit() {
	if that.Save != nil {
		//提交事务
		that.Save.Commit()
	}
}
func (that *BaseWsController) Close() {
	if that.Save != nil {
		//关闭事务
		that.Save.Close()
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
func (that *BaseWsController) Start(serviceFunc func() string) (result string) {
	result = "{\"code\":1,\"message\":\"系统错误！\",\"data\":\"\"}"
	that.Begin()
	defer that.handleException(&result)
	defer that.Close()
	result = serviceFunc()
	that.Commit()
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
	temp := "{data: \"\", code: 1, message: \"服务器出现异常，请稍后再试！\"}"
	return &temp
}

/**
处理异常信息
*/
func (that *BaseWsController) handleException(result *string) {
	if err := recover(); err != nil {
		if that.Save != nil {
			//事务回滚
			that.Save.Rollback()
		}
		//异常处理
		result = that.exceptionRecover(err)
	} else {
		temp := "{data: \"handleException\", code: 0, message: \"\"}"
		result = &temp
	}
}
