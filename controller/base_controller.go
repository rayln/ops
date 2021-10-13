package controller

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/entity"
	"github.com/rayln/ops/util"
	"runtime"
)

type BaseController struct {
	entity.BaseEntity
	Request context.Context
}

/**
开启事务
*/
func (that *BaseController) Begin() {
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
func (that *BaseController) Commit() {
	if that.Save != nil {
		//提交事务
		that.Save.Commit()
	}
}
func (that *BaseController) Close() {
	if that.Save != nil {
		//关闭事务
		that.Save.Close()
		that.Save = nil
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
func (that *BaseController) Start(request iris.Context, serviceFunc func() mvc.Result) mvc.Result {
	that.Begin()
	defer that.Close()
	defer that.handleException(request)
	result := serviceFunc()
	that.Commit()
	return result
}

func (that *BaseController) BeforeActivation(a mvc.BeforeActivation) {
	//拦截链接，主要用于Session，开启一个Session，并且捕获异常进行回滚
	a.Router().Use(func(context context.Context) {
		context.Next()
	})
	a.Router().Done(func(context context.Context) {
		//跳转结束
		context.Next()
	})
}
func (that *BaseController) AfterActivation(a mvc.AfterActivation) {

}

/**
* @Description: 错误信息处理
 */
func (that *BaseController) exceptionRecover(ctx iris.Context, err interface{}) {
	if ctx.IsStopped() {
		return
	}

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
	logMessage := fmt.Sprintf("从错误中回复：('%s')\n", ctx.HandlerName())
	logMessage += errMsg + "\n"
	logMessage += fmt.Sprintf("\n%s", stacktrace)
	that.Logger.Error(logMessage)
	// 打印错误日志
	// 返回错误信息
	ctx.JSON(util.Result{Data: "", Code: util.ERROR_CODE, Message: "服务器出现异常，请稍后再试！"})
	ctx.StatusCode(100)
	// 停止跳转
	//ctx.StopExecution()
}

/**
处理异常信息
*/
func (that *BaseController) handleException(request context.Context) {
	if err := recover(); err != nil {
		if that.Save != nil {
			//事务回滚
			that.Save.Rollback()
		}
		//异常处理
		that.exceptionRecover(request, err)
		return
	}
}
