package job

import (
	"context"
	"github.com/pkg/errors"
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
	strategy *CompositeStrategy
	started  uint32
	cancel   func()
	done     chan struct{}
}

func New(payload Payload, strategies ...Strategy) (job *Job) {
	job = &Job{
		payload:  payload,
		strategy: Compose(strategies...),
		started:  no,
		done:     make(chan struct{}),
	}
	return
}

func (p *Job) Start() {
	p.StartContext(context.Background())
}

var alreadyStartedError = errors.New("already started")

func (p *Job) StartContext(ctx context.Context) {
	if !atomic.CompareAndSwapUint32(&p.started, no, yes) {
		panic(alreadyStartedError)
	}
	ctx, p.cancel = context.WithCancel(ctx)

	go func() {
		p.run(ctx)
		close(p.done)
	}()
	return
}

func (p *Job) run(ctx context.Context) {
	lastTickTime := time.Now()
	nextTickTime, ok := p.strategy.Tick(lastTickTime)
	if !ok {
		return
	}
	timer := time.NewTimer(time.Until(nextTickTime))
	defer timer.Stop()
	if !wait(ctx, timer.C) {
		return
	}

	for {
		p.payload(ctx)

		lastTickTime = nextTickTime
		nextTickTime, ok = p.strategy.Tick(lastTickTime)
		if !ok {
			return
		}
		timer.Reset(time.Until(nextTickTime))
		if !wait(ctx, timer.C) {
			return
		}
	}
}

func (p *Job) Done() (done <-chan struct{}) {
	done = p.done
	return
}

var notStartedError = errors.New("not started")

func (p *Job) Stop() {
	p.StopContext(context.Background())
}

func (p *Job) StopContext(ctx context.Context) {
	if atomic.LoadUint32(&p.started) == no {
		panic(notStartedError)
	}

	p.cancel()
	wait(ctx, p.done)
	return
}

func wait[T any](ctx context.Context, ch <-chan T) (ok bool) {
	select {
	case <-ctx.Done():
		break
	case <-ch:
		ok = true
		break
	}
	return
}
