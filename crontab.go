package crontab

import (
	"log"
	"time"
)

//定时任务结构体
type Job struct {
	Interval     int
	IntervalTime time.Duration
	Ticker       *time.Ticker
	Timer        *time.Timer
	Handler      func()
	StopChan     chan bool
}

//运行多次任务
func (job *Job) RunTicker() {
	job.Ticker = time.NewTicker(job.IntervalTime)
	defer job.Ticker.Stop()

	for {
		select {
		case <-job.Ticker.C:
			job.Handler()
		case stop := <-job.StopChan:
			if stop {
				log.Println("这个定时器停了")
				return
			}
		}
	}
}

//运行单次任务
func (job *Job) RunTimer() {
	job.Timer = time.NewTimer(job.IntervalTime)
	defer job.Timer.Stop()

	for range job.Timer.C {
		job.Handler()
	}
}

//运行任务自身
func (job *Job) Run() {
	job.StopChan = make(chan bool)

	switch job.Interval {
	case TypeTicker:
		go job.RunTicker()
	case TypeTimer:
		go job.RunTimer()
	}
}

//停止当前任务
func (job *Job) Stop() {
	//time.Sleep(time.Microsecond * 1)  //sleep一段时间会更好的退出
	job.StopChan <- true
	close(job.StopChan)
}

//----------------------------

//封装定时执行程序结构体
type Crontab struct {
	list []*Job
	init bool
}

//添加一个定时程序
func (c *Crontab) Add(intervalType int, intervalTime time.Duration, callback func()) (job *Job) {
	job = &Job{
		Interval:     intervalType,
		IntervalTime: intervalTime,
		Handler:      callback,
		StopChan:     make(chan bool),
		Ticker:       nil,
		Timer:        nil,
	}

	switch job.Interval {
	case TypeTimer:
		job.Timer = time.NewTimer(intervalTime)
	case TypeTicker:
		job.Ticker = time.NewTicker(intervalTime)
	}

	c.list = append(c.list, job)
	return
}

//----------------------------

//定义单次定时器类型
const TypeTimer = 0

//定义多次定时器类型
const TypeTicker = 1

var (
	cron = new(Crontab)
)

//创建一个定时器
func New(intervalType int, intervalTime time.Duration, callback func()) {
	job := cron.Add(intervalType, intervalTime, callback)

	// 启动状态则直接启动该定时器任务
	if cron.init {
		job.Run()
	}
}

//创建一个单次定时器
func NewTimer(intervalTime time.Duration, callback func()) {
	New(TypeTimer, intervalTime, callback)
}

//创建一个多次定时器
func NewTicker(intervalTime time.Duration, callback func()) {
	New(TypeTicker, intervalTime, callback)
}

//初始化已创建的定时器
func Init() {
	cron.init = true

	for _, job := range cron.list {
		job.Run()
	}
}

//停止某一定时器或所有定时器
func Stop(cronId string) {
	//Todo: 能够停止某一定时器或所有定时器
}
