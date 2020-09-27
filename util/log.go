package util

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/middleware/logger"
	"io"
	"os"
	"time"
)

type Log struct {
	fileTime time.Time
	logFile  *os.File
}

func (that *Log) Init(app *iris.Application) {
	that.createfile(app)
	that.createHandler(app)
	that.checkfile(app)
	return
}
func (that *Log) createfile(app *iris.Application) {
	that.logFile = that.newLogFile()
	app.Logger().SetOutput(io.MultiWriter(that.logFile, os.Stdout))
	app.Logger().SetLevel("debug")
	that.fileTime = time.Now()
}

//创建一个中间键，获取一个请求的时长
func (that *Log) createHandler(app *iris.Application) {
	var h iris.Handler
	c := logger.Config{
		Status:  true,
		IP:      true,
		Method:  true,
		Path:    true,
		Columns: true,
	}
	c.LogFunc = func(now time.Time, latency time.Duration, status, ip, method, path string, message interface{}, headerMessage interface{}) {
		that.logFile.Write([]byte(fmt.Sprintln(now.Format("[HTTP] 2006/01/02 15:04"), " 耗时:", latency, " 状态:", status, " IP地址:", ip, " 请求:", method, " 路径:", path)))
	}
	h = logger.New(c)
	app.Use(h)
}
func (that *Log) checkfile(app *iris.Application) {
	//每60秒检测一次，是否需要新开一个Log文件。
	new(Schedule).Run(func() {
		now := time.Now()
		if that.fileTime.Day() != now.Day() {
			//创建一个新文件（判断日期不同）
			that.createfile(app)
			that.fileTime = now
		}
	}, time.Second*60)
}

func (that *Log) newLogFile() *os.File {
	path, _ := os.Getwd()
	configDir := path + "/file/logfile"
	os.MkdirAll(configDir, os.ModeDir)

	filename := that.todayFilename()
	//打开一个输出文件，如果重新启动服务器，它将追加到今天的文件中
	f, err := os.OpenFile(configDir+"/"+filename, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}
	return f
}
func (that *Log) todayFilename() string {
	today := time.Now().Format("2006_01_02")
	return today + ".log"
}
