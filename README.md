# crontab
golang封装定时器

## 依赖要求
没有

## 安装
使用`go`命令获取类库

```bash
go get github.com/flxxyz/crontab
```

## 例子
```go
package main

import (
    "github.com/flxxyz/crontab"
    "log"
    "time"
)

func init() {
    //运行定时器，可随意放置代码位置
    crontab.Init()
}

func main() {
    //创建一个单次定时器
    crontab.NewTimer(time.Second*5, func() {
        log.Println("这是一个单次定时器")
    })

    //创建一个多次定时器
    crontab.NewTicker(time.Second*1, func() {
        log.Println("这是一个多次定时器")
    })
}
```

## 文档
[文档点这里](http://godoc.org/github.com/flxxyz/crontab)

## 版权
crontab包在MIT License下发布。有关详细信息，请参阅LICENSE。