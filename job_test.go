package job

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func Test_OnStartJobWithoutDelay_ShouldInvokePayloadFunctionImmediately(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime
	}))

	go job.Start()
	defer job.Stop()
	time.Sleep(time.Second)
	assert.NotEqual(t, uint32(0), atomic.LoadUint32(&counter))
}

func Test_OnStartJobWithDelay_ShouldInvokePayloadFunctionWithDelay(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(2 * time.Second)
	}))

	go job.Start()
	defer job.Stop()
	time.Sleep(time.Second)
	assert.Equal(t, uint32(0), atomic.LoadUint32(&counter))
	time.Sleep(2 * time.Second)
	assert.Equal(t, uint32(1), atomic.LoadUint32(&counter))
}

func Test_OnStartJobWithPeriod_ShouldInvokePayloadFunctionPeriodically(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(3 * time.Second)
	}))

	go job.Start()
	defer job.Stop()
	time.Sleep(time.Second)
	for i := 0; i < 4; i++ {
		assert.Equal(t, uint32(i), atomic.LoadUint32(&counter))
		time.Sleep(3 * time.Second)
	}
}

func Test_OnStartJobAndCancelContext_ShouldStopJob(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(2 * time.Second)
	}))

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	go job.StartContext(ctx)
	<-job.Done()
	cancel()
	assert.Equal(t, uint32(2), atomic.LoadUint32(&counter))
}

func Test_OnStopJob_ShouldCancelExecutionContext(t *testing.T) {
	var counter uint32
	job := New(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			atomic.AddUint32(&counter, 1)
		case <-time.After(2 * time.Second):
			break
		}
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime
	}))

	go job.Start()
	time.Sleep(time.Second)
	job.Stop()
	assert.Equal(t, uint32(1), atomic.LoadUint32(&counter))
}

func Test_OnStopJobBetweenTicks_ShouldStopJobImmediately(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(2 * time.Second)
	}))

	go job.Start()
	time.Sleep(3 * time.Second)
	job.Stop()
	assert.Equal(t, uint32(1), atomic.LoadUint32(&counter))
}

func Test_OnStopJobDuringTick_ShouldWaitForJobGracefulShutdown(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		time.Sleep(2 * time.Second)
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime
	}))

	go job.Start()
	time.Sleep(time.Second)
	job.Stop()
	assert.Equal(t, uint32(1), atomic.LoadUint32(&counter))
}

func Test_OnStopJobWithCancelledContext_ShouldStopWaitingForShutdown(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		time.Sleep(3 * time.Second)
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime
	}))

	go job.Start()
	time.Sleep(time.Second)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	job.StopContext(ctx)
	cancel()
	assert.Equal(t, uint32(0), atomic.LoadUint32(&counter))
}

func Test_OnStartGoneJob_ShouldExitImmediately(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime
	}))

	job.Stop()
	job.Start()
	assert.Equal(t, uint32(0), atomic.LoadUint32(&counter))
}

func Test_OnJobFinish_ShouldSendDoneSignal(t *testing.T) {
	var counter uint32
	job := New(func(_ context.Context) {
		atomic.AddUint32(&counter, 1)
	}, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime
	}))

	go job.Start()
	time.Sleep(time.Second)
	job.Stop()
	<-job.Done()
	assert.NotEqual(t, uint32(0), atomic.LoadUint32(&counter))
}
