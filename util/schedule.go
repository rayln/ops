package util

import "time"

type Schedule struct {
	isEnd bool
}

/**
定时器，根据传入的时间，进行无限循环执行任务
*/
func (that *Schedule) Run(callfunc func(), duration time.Duration) {
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
