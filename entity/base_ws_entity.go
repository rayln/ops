package entity

import (
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/rayln/ops/util"
	"xorm.io/xorm"
)

type BaseWsEntity struct {
	//SessionTmp *xorm.Session
	Engine      *xorm.EngineGroup
	Logger      *golog.Logger
	Redis       *util.Redis
	App         *iris.Application
	EngineOther []*xorm.EngineGroup
}
