package main

import (
	"context"
	"time"
)

type Job struct {
	cancel context.CancelFunc

	postHooks []func()

	totalUnits int64
	unitTime   time.Duration
}

func NewJob(totalUnits int64, unitTime time.Duration) *Job {
	return &Job{
		totalUnits: totalUnits,
		unitTime:   unitTime,
		cancel:     nil,
	}
}

func (j *Job) Total() int64 {
	return j.totalUnits
}

func (j *Job) OnUnitCompletion(hook func()) {
	j.postHooks = append(j.postHooks, hook)
}

func (j *Job) StartAsync(ctx context.Context) {
	_ctx, cancel := context.WithCancel(ctx)

	j.cancel = cancel

	go func() {
		for i := int64(0); i < j.totalUnits; i++ {
			select {
			case <-_ctx.Done():
				return
			default:
				time.Sleep(j.unitTime)
				for _, hook := range j.postHooks {
					hook()
				}
			}
		}
	}()
}

func (j *Job) Stop() {
	if j.cancel != nil {
		j.cancel()
	}
}
