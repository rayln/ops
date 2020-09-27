package config

import (
	"github.com/go-xorm/xorm"
	"github.com/rayln/ops/entity"
)

func InitDatabase(engine *xorm.Engine) {
	//初始化数据库表
	err := engine.Sync(
		new(entity.FarmTestInfo),
	)
	if err != nil {
		panic(err.Error())
	}
	//engine.ImportFile("./com/sql/farm_constant_info.sql")
}
