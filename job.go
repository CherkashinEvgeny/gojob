package job

import (
	"context"
	"github.com/pkg/errors"
	"sync"
)

type Payload func(ctx context.Context)

type Status uint

const (
	Ready Status = iota
	Running
	Stopped
)

type Job struct {
	payload Payload
	status  Status
	cancel  func()
	done    chan struct{}
	lock    *sync.RWMutex
}

var goneDoneChan = make(chan struct{})

func New(payload Payload) (job *Job) {
	job = &Job{
		payload: payload,
		status:  Ready,
		done:    goneDoneChan,
		lock:    &sync.RWMutex{},
	}
	return
}

func (p *Job) Status() (status Status) {
	p.lock.RLock()
	status = p.status
	p.lock.RUnlock()
	return
}

func (p *Job) Start(strategies ...Strategy) (err error) {
	err = p.StartContext(context.Background(), strategies...)
	return
}

var NotReadyError = errors.New("job is not ready to start")

func (p *Job) StartContext(ctx context.Context, strategies ...Strategy) (err error) {
	p.lock.Lock()
	if p.status != Ready {
		p.lock.Unlock()
		err = NotReadyError
		return
	}
	ctx, p.cancel = context.WithCancel(ctx)
	p.done = make(chan struct{})
	p.status = Running
	p.lock.Unlock()

	go func() {
		p.run(ctx, strategies...)
		p.lock.Lock()
		close(p.done)
		p.status = Ready
		p.lock.Unlock()
	}()
	return
}

func (p *Job) run(ctx context.Context, strategies ...Strategy) {
	strategy := Compose(strategies...)
	defer strategy.Reset()
	for !isContextCancelled(ctx) && strategy.Tick(ctx) {
		p.payload(ctx)
	}
}

func isContextCancelled(ctx context.Context) (cancelled bool) {
	select {
	case <-ctx.Done():
		cancelled = true
		break
	default:
		break
	}
	return
}

func (p *Job) Done() (done <-chan struct{}) {
	p.lock.RLock()
	done = p.done
	p.lock.RUnlock()
	return
}

func (p *Job) Stop() (err error) {
	err = p.StopContext(context.Background())
	return
}

var NotRunningError = errors.New("not running")

func (p *Job) StopContext(ctx context.Context) (err error) {
	p.lock.Lock()
	if p.status != Running {
		p.lock.RUnlock()
		err = NotRunningError
		return
	}
	p.status = Stopped
	done := p.done
	cancel := p.cancel
	p.lock.Unlock()

	cancel()
	select {
	case <-ctx.Done():
		err = ctx.Err()
		break
	case <-done:
		break
	}
	return
}
