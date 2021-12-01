package entity

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/rayln/ops/util"
	"xorm.io/xorm"
)

type BaseEntity struct {
	//SessionTmp *xorm.Session
	Engine      *xorm.EngineGroup
	Logger      *golog.Logger
	Save        *xorm.Session
	Load        *xorm.Engine
	Redis       *util.Redis
	App         *iris.Application
	EngineOther []*xorm.EngineGroup
	SaveOther   []*xorm.EngineGroup
	LoadOther   []*xorm.Engine
}

var IsClose = false
