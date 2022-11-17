# ziva

一个go实现的多任务、多线程爬虫库


## 一个demo

```go
package main

import (
	"fmt"
	"github.com/zituocn/ziva"
	"net/http"
)

func main() {
	job := ziva.NewJob("article", ziva.Options{
		CreateQueue: func() ziva.TodoQueue {
			ids := []int{3263, 3262, 3261, 3260, 3259}
			queue := ziva.NewMemQueue()
			header := &http.Header{}
			header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
			for _, item := range ids {
				queue.Add(&ziva.Task{
					Url:    fmt.Sprintf("%s%d", "https://22v.net/article/", item),
					Method: "GET",
					Header: header,
				})
			}
			return queue
		},
		SucceedFunc: func(ctx *ziva.Context) {
			fmt.Println("成功的回调")
			fmt.Println("返回信息 :", ctx.Response.Status)
		},
		FailedFunc: func(ctx *ziva.Context) {
			fmt.Println("失败的回调")
			fmt.Println("返回状态 :", ctx.Response.StatusCode)
		},
		SheepTime: 3000,
		Num:       1,
	})

	job.Do()
}

```

## 使用redis queue

```go
package main

import (
	"fmt"
	"net/http"

	"github.com/zituocn/ziva"
	"github.com/zituocn/ziva/goredis"
)

func main() {
	job := ziva.NewJob("article", ziva.Options{
		CreateQueue: func() ziva.TodoQueue {
			ids := []int{3263, 3262, 3261, 3260, 3259}
			queue := ziva.NewRedisQueue("article", &goredis.RedisConfig{
				Name:     "127.0.0.1",
				Host:     "127.0.0.1",
				Port:     6379,
				Password: "123456",
				DB:       1,
			})
			header := &http.Header{}
			header.Set("User-Agent", "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_15_7) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/100.0.4896.127 Safari/537.36")
			for _, item := range ids {
				queue.Add(&ziva.Task{
					Url:    fmt.Sprintf("%s%d", "https://22v.net/article/", item),
					Method: "GET",
					Header: header,
				})
			}
			return queue
		},
		SucceedFunc: func(ctx *ziva.Context) {
			fmt.Println("成功的回调")
			fmt.Println("返回信息 :", ctx.Response.Status)
		},
		FailedFunc: func(ctx *ziva.Context) {
			fmt.Println("失败的回调")
			fmt.Println("返回状态 :", ctx.Response.StatusCode)
		},
		SheepTime: 3000,
		IsDebug:   true,
		Num:       1,
	})

	job.Do()
}

```