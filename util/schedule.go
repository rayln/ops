package util

import "time"

type Schedule struct {
}

/**
定时器，根据传入的时间，进行无限循环执行任务
*/
func (that *Schedule) Run(callfunc func(), duration time.Duration) {
	go func() {
		for {
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
	go func() {
	Loop:
		for {
			tiker := time.NewTicker(duration)
			if callfunc() {
				break Loop
			}
			<-tiker.C
		}
	}()
}
