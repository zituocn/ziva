package ziva

import (
	"github.com/zituocn/ziva/logx"
	"sync"
)

type Job struct {
	name string

	options Options
}

func NewJob(name string, options Options) *Job {
	if options.Num < 1 {
		options.Num = 1
	}
	j := &Job{
		name:    name,
		options: options,
	}
	if j.options.CreateQueue != nil {
		j.options.Queue = j.options.CreateQueue()
	}
	return j
}

func (j *Job) Do() {
	logx.Infof("[%s] start job -> Goroutines: %d ", j.name, j.options.Num)
	var wg sync.WaitGroup
	for n := 0; n < j.options.Num; n++ {
		wg.Add(1)
		go func(i int) {
			logx.Infof("start task %d", i+1)
			defer wg.Done()
			for {
				if j.options.Queue.IsEmpty() {
					break
				}
				task := j.options.Queue.Pop()
				if task != nil {
					ctx := DoTask(task)
					ctx.Options = j.options
					ctx.Do()
				}
			}
		}(n)
	}
	wg.Wait()
	logx.Info("job done")
}
