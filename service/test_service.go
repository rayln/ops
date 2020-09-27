package service

import (
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/dao"
	"github.com/rayln/ops/entity"
	"github.com/rayln/ops/util"
)

type TestService struct {
	BaseService
	TestDao dao.TestDao
	//定义一个返回值
	Result util.Result
}

/***************************逻辑处理开始*******************************/
func (that *TestService) Test(base *entity.BaseEntity) mvc.Result {
	/**===============逻辑开始===========================**/
	//保存记录
	that.TestDao.Save(base)
	//查询记录
	var farmTestInfo = that.TestDao.Query(base)
	base.Logger.Info("farmTestInfo: ", farmTestInfo)
	//redis操作
	that.TestDao.RedisTest(base)
	/**===============逻辑结束===========================**/
	//返回值
	that.Result.Data = "插入成功！！"
	return that.Result.Response()
}
