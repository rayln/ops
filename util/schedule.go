package util

import (
	"fmt"
	"runtime"
	"time"
)

type Schedule struct {
	isEnd bool
}

/**
定时器，根据传入的时间，进行无限循环执行任务
*/
func (that *Schedule) Run(callfunc func(), duration time.Duration) {
	defer that.handleException()
	that.isEnd = false
	go func() {
		tiker := time.NewTicker(duration)
		defer tiker.Stop()
		for {
			callfunc()
			<-tiker.C
			if that.isEnd {
				break
			}
		}
	}()
}

/**
返回true则跳出循环
*/
func (that *Schedule) RunToBreak(callfunc func() bool, duration time.Duration) {
	defer that.handleException()
	that.isEnd = false
	go func() {
		tiker := time.NewTicker(duration)
		defer tiker.Stop()
		for {
			<-tiker.C
			if callfunc() {
				break
			}
			if that.isEnd {
				break
			}
		}
	}()
}

func (that *Schedule) Delay(callfunc func(), duration time.Duration) {
	defer that.handleException()
	that.isEnd = false
	go func() {
		tiker := time.NewTicker(duration)
		defer tiker.Stop()
		for {
			<-tiker.C
			callfunc()
			break
		}
	}()
}

/**
移除所有定时任务。无法在下达的时候及时停止。但是会在下一时刻停止
*/
func (that *Schedule) RemoveAllSchedule() {
	that.isEnd = true
}

func (that *Schedule) handleException() {

	if err := recover(); err != nil {
		that.exceptionRecover(err)
		return
	}
}

func (that *Schedule) exceptionRecover(err interface{}) {

	var stacktrace string
	for i := 1; ; i++ {
		_, f, l, got := runtime.Caller(i)
		if !got {
			break
		}
		stacktrace += fmt.Sprintf("%s:%d\n", f, l)
	}

	errMsg := fmt.Sprintf("错误信息: %s", err)
	logMessage := fmt.Sprintf("从错误中回复：\n")
	logMessage += errMsg + "\n"
	logMessage += fmt.Sprintf("\n%s", stacktrace)
	fmt.Println(logMessage)
}
