package job

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

type Payload func(ctx context.Context)

const (
	no = iota
	yes
)

type Job struct {
	payload  Payload
	strategy Strategy
	started  uint32
	used     uint32
	stopped  uint32
	cancel   func()
	done     chan struct{}
}

func New(payload Payload, strategy Strategy) (job *Job) {
	return &Job{
		payload:  payload,
		strategy: strategy,
		started:  no,
		used:     no,
		stopped:  no,
		cancel:   nil,
		done:     make(chan struct{}),
	}
}

func (j *Job) Start() {
	j.StartContext(context.Background())
}

var alreadyStartedError = errors.New("already started")

func (j *Job) StartContext(ctx context.Context) {
	if !atomic.CompareAndSwapUint32(&j.started, no, yes) {
		panic(alreadyStartedError)
	}
	ctx, j.cancel = context.WithCancel(ctx)
	if !atomic.CompareAndSwapUint32(&j.used, no, yes) {
		// stop command before initialization
		// releases context resources
		j.cancel()
		return
	}
	defer func() {
		j.cancel()
		close(j.done)
	}()
	j.run(ctx)
}

func (j *Job) run(ctx context.Context) {
	lastTickTime := time.Now()
	nextTickTime := j.strategy.Tick(lastTickTime)
	timer := time.NewTimer(time.Until(nextTickTime))
	defer timer.Stop()
	for {
		if !waitForTimerSignal(ctx, timer.C) {
			return
		}
		j.payload(ctx)
		lastTickTime = nextTickTime
		nextTickTime = j.strategy.Tick(lastTickTime)
		timer.Reset(time.Until(nextTickTime))
	}
}

func waitForTimerSignal(ctx context.Context, ch <-chan time.Time) (ok bool) {
	select {
	case <-ctx.Done():
		return false
	case <-ch:
		return true
	}
}

func (j *Job) Done() (done <-chan struct{}) {
	return j.done
}

func (j *Job) Stop() {
	j.StopContext(context.Background())
}

func (j *Job) StopContext(ctx context.Context) {
	if !atomic.CompareAndSwapUint32(&j.stopped, no, yes) {
		j.waitForGracefulShutdown(ctx)
		return
	}
	if atomic.CompareAndSwapUint32(&j.used, no, yes) {
		// stop command before initialization
		close(j.done)
		return
	}
	j.cancel()
	j.waitForGracefulShutdown(ctx)
}

func (j *Job) waitForGracefulShutdown(ctx context.Context) (ok bool) {
	select {
	case <-ctx.Done():
		return false
	case <-j.done:
		return true
	}
}
