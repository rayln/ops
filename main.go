package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/mvc"
	"github.com/rayln/ops/config"
	"github.com/rayln/ops/entity"
	"github.com/rayln/ops/intercepter"
	"github.com/rayln/ops/util"
	"time"
	"xorm.io/core"
)

func main() {
	app := newApp()
	//设置配置文件
	appconfig := new(config.AppConfig).Init()
	dataconfig := new(config.DatabaseConfig).Init()
	//设置log日志
	new(util.Log).Init(app)
	//初始化服务器
	initserver(app)
	//设置路由
	handle(app, appconfig, dataconfig)

	app.Run(
		iris.Addr(appconfig.Port),
		iris.WithCharset("UTF-8"),
		iris.WithoutServerError(iris.ErrServerClosed), //无服务器错误提示
		iris.WithOptimizations,                        //对json数据序列化更快配置
		iris.WithoutBodyConsumptionOnUnmarshal,
	)
}

/**
APP的创建
*/
func newApp() *iris.Application {
	app := iris.New()
	/*
	* requestPath是请求路径。systemPath是实际路径
	* 比如 localhost:8080/manager/static/rayln.png。文件路径是/static/rayln.png
	* 那么requestPath是/manager/static/rayln.png，systemPath是/static/rayln.png
	* app.StatWeb("/manager/static","./static")
	 */
	//app.StaticWeb("/static","./static")
	//注册html页面
	app.RegisterView(iris.HTML("./static", ".html"))
	//欢迎页面，如果用户不输入index.html页面，默认访问index.html
	app.Get("/", func(context context.Context) {
		context.View("index.html")
	})

	return app
}

/**
config设置
*/
func initserver(app *iris.Application) {
	//配置字符编码
	app.Configure(iris.WithConfiguration(iris.Configuration{
		Charset: "UTF-8",
	}))
	//404错误配置
	app.OnErrorCode(iris.StatusNotFound, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusNotFound,
			"msg":    "not found",
			"data":   iris.Map{},
		})
	})
	//500系统错误
	app.OnErrorCode(iris.StatusInternalServerError, func(context context.Context) {
		context.JSON(iris.Map{
			"errmsg": iris.StatusInternalServerError,
			"msg":    "interal error",
			"data":   iris.Map{},
		})
	})
}

/**
设置路由
*/
func handle(app *iris.Application, appConfig *config.AppConfig, databaseConfig *config.DatabaseConfig) {
	engine := InitDatabase(databaseConfig)
	redis := new(util.Redis).Init(databaseConfig.Redis, databaseConfig.RedisPwd, databaseConfig)
	application := mvc.New(app.Party("/"))
	//重要
	//所有选项都可以用Force填充：true，所有的都会很好的兼容
	application.Router.SetExecutionRules(iris.ExecutionRules{
		// Begin: <- from `Use[all]` 到`Handle[last]` 程序执行顺序，执行all，即使缺少`ctx.Next()`也执行all。
		// Main: <- all `Handle` 执行顺序，执行所有>> >>。
		// < - 从`Handle [last]`到`Done [all]`程序执行顺序，执行全部>> >>。
		Done: iris.ExecutionOptions{Force: true},
	})
	application.Router.Use(func(context context.Context) {
		//拦截器的使用
		if appConfig.IsSign {
			var pass = new(intercepter.Intercepter).Init().Handle(context)
			if !pass {
				context.JSON(util.Result{Data: "", Code: util.ERROR_INTERCEPT_MD5, Message: "system error"})
				return
			}
		} else {
			if entity.IsClose {
				context.JSON(util.Result{Data: "", Code: util.ERROR_INTERCEPT_MD5, Message: "system is closing"})
				return
			}
		}
		//拦截器对业务的处理
		context.Next()
	})
	application.Router.Done(func(ctx context.Context) {
		//主逻辑完成后到这里做最后处理（无论业务逻辑是否抛出异常，都在到这里）
	})
	//绑定Controller的参数，暂时为null
	application.Register(app.Logger(), engine, redis)
	config.InitController(application, app, engine)
}

/**
初始化数据库
*/
func InitDatabase(databaseConfig *config.DatabaseConfig) *xorm.EngineGroup {
	conns := []string{
		fmt.Sprintf("%s", databaseConfig.Master), //腾讯云服务器读写实例
		fmt.Sprintf("%s", databaseConfig.Select), //从库
	}
	// 第三个参数是负载策略，当前LeastConnPolicy是最小连接数访问负载策略
	engine, err := xorm.NewEngineGroup("mysql", conns, xorm.LeastConnPolicy())
	if err != nil {
		panic(err.Error())
	}
	//数据库设置
	//显示SQL语句
	engine.ShowSQL(databaseConfig.ShowSQL)
	//设置日志级别
	if databaseConfig.LogLevel == 1 {
		engine.Logger().SetLevel(core.LOG_DEBUG)
	} else if databaseConfig.LogLevel == 2 {
		engine.Logger().SetLevel(core.LOG_INFO)
	} else {
		engine.Logger().SetLevel(core.LOG_INFO)
	}
	//设置连接池空闲数大小
	engine.SetMaxIdleConns(databaseConfig.MaxIdleConns)
	//设置连接池最大连接数
	engine.SetMaxOpenConns(databaseConfig.MaxOpenConns)
	//设置连接池可以使用的最长有效时间
	//engine.SetConnMaxLifetime(50000)
	//设置映射规则（驼峰规则）
	engine.SetMapper(core.SnakeMapper{})
	//配置表结构
	config.InitDatabase(engine.Master())
	//缓存处理
	//xrc.NewRedisCacher("localhost:6379", "", xrc.DEFAULT_EXPIRATION, engine.Logger())
	cache := xorm.NewLRUCacher2(xorm.NewMemoryStore(), time.Duration(3600)*time.Second, 10000)
	engine.SetDefaultCacher(cache)
	/*f,_ := os.Create("sql.log")
	xorm.Engine.Logger = xorm.NewSimpleLogger(f)*/
	return engine
}
