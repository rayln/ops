package entity

import (
	"github.com/go-xorm/xorm"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/rayln/ops/util"
)

type BaseEntity struct {
	//SessionTmp *xorm.Session
	Engine *xorm.EngineGroup
	Logger *golog.Logger
	Save   *xorm.Session
	Load   *xorm.Engine
	Redis  *util.Redis
	App    *iris.Application
}

var IsClose = false
