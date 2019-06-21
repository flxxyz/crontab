package crontab

import (
    "log"
    "time"
)

type Job struct {
    Interval     int
    IntervalTime time.Duration
    Ticker       *time.Ticker
    Timer        *time.Timer
    Handler      func()
    StopChan     chan bool
}

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

func (job *Job) RunTimer() {
    job.Timer = time.NewTimer(job.IntervalTime)
    defer job.Timer.Stop()

    for range job.Timer.C {
        job.Handler()
    }
}

func (job *Job) Run() {
    job.StopChan = make(chan bool)

    switch job.Interval {
    case TypeTicker:
        go job.RunTicker()
    case TypeTimer:
        go job.RunTimer()
    }
}

func (job *Job) Stop() {
    //time.Sleep(time.Microsecond * 1)  //sleep一段时间会更好的退出
    job.StopChan <- true
    close(job.StopChan)
}

//----------------------------

type Crontab struct {
    list []*Job
    run  bool
}

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

const TypeTimer = 0
const TypeTicker = 1

var (
    crontab = new(Crontab)
)

func New(intervalType int, intervalTime time.Duration, callback func()) {
    job := crontab.Add(intervalType, intervalTime, callback)

    // 启动状态则直接启动该定时器任务
    if crontab.run {
        job.Run()
    }
}

func NewTimer(intervalTime time.Duration, callback func()) {
    New(TypeTimer, intervalTime, callback)
}

func NewTicker(intervalTime time.Duration, callback func()) {
    New(TypeTicker, intervalTime, callback)
}

func Run() {
    crontab.run = true

    for _, c := range crontab.list {
        c.Run()
    }
}

func Stop(cronId string) {
    //Todo: 能够停止某一定时器或所有定时器
}
