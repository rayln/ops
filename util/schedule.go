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
		for {
			if that.isEnd {
				break
			}
			tiker := time.NewTicker(duration)
			callfunc()
			<-tiker.C
		}
	}()
}

/**
返回true则跳出循环
*/
func (that *Schedule) RunToBreak(callfunc func() bool, duration time.Duration) {
	that.isEnd = false
	go func() {
	Loop:
		for {
			if that.isEnd {
				break
			}
			tiker := time.NewTicker(duration)
			<-tiker.C
			if callfunc() {
				break Loop
			}

		}
	}()
}

func (that *Schedule) Delay(callfunc func(), duration time.Duration) {
	that.isEnd = false
	go func() {
	Loop:
		for {
			if that.isEnd {
				break
			}
			tiker := time.NewTicker(duration)
			<-tiker.C
			callfunc()
			break Loop
		}
	}()
}

func (that *Schedule) RemoveAllSchedule() {
	that.isEnd = true
}
