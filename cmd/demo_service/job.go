package main

import (
	"context"
	"time"

	"github.com/IceWhaleTech/CasaOS-JobManagement/codegen"
)

type Task struct {
	cancel           context.CancelFunc
	onUnitCompletion []func()
	jobID            codegen.JobID
	totalUnits       int64
	unitTime         time.Duration
}

func NewTask(totalUnits int64, unitTime time.Duration) *Task {
	return &Task{
		cancel:           nil,
		jobID:            -1,
		onUnitCompletion: make([]func(), 0),
		totalUnits:       totalUnits,
		unitTime:         unitTime,
	}
}

func (t *Task) StartAsync(ctx context.Context, onStart func(_t *Task)) {
	_ctx, cancel := context.WithCancel(ctx)

	t.cancel = cancel

	go func() {
		onStart(t)

		for i := int64(0); i < t.totalUnits; i++ {
			select {
			case <-_ctx.Done():
				return
			default:
				time.Sleep(t.unitTime)
				for _, hook := range t.onUnitCompletion {
					hook()
				}
			}
		}
	}()
}

func (t *Task) Stop() {
	if t.cancel == nil {
		return
	}

	t.cancel()
}
