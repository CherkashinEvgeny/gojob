package job

import (
	"context"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_OnStartJobWithoutDelay_ShouldInvokePayloadFunctionImmediately(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		done = true
	}, NoDelay())
	job.Start()

	// wait for start
	time.Sleep(time.Second)

	assert.True(t, done)
}

func Test_OnStartJobWithDelay_ShouldInvokePayloadFunctionWithDelay(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		done = true
	}, Delay(2*time.Second))
	job.Start()

	// wait for start
	time.Sleep(time.Second)
	assert.False(t, done)

	//wait for delay
	time.Sleep(2 * time.Second)
	assert.True(t, done)
}

func Test_OnStartJobWithPeriod_ShouldInvokePayloadFunctionPeriodically(t *testing.T) {
	counter := 0
	job := New(func(_ context.Context) {
		counter++
	}, Interval(3*time.Second))
	job.Start()
	defer job.Stop()

	// wait for start
	time.Sleep(time.Second)

	for i := 0; i < 4; i++ {
		assert.Equal(t, i, counter)
		time.Sleep(3 * time.Second)
	}
}

func Test_OnStartJobAndCancelContext_ShouldStopJob(t *testing.T) {
	counter := 0
	job := New(func(_ context.Context) {
		counter++
	}, Interval(2*time.Second))
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	job.StartContext(ctx)
	<-job.Done()
	cancel()
	assert.Equal(t, 2, counter)
}

func Test_OnStopJob_ShouldCancelExecutionContext(t *testing.T) {
	canceled := false
	job := New(func(ctx context.Context) {
		select {
		case <-ctx.Done():
			canceled = true
		case <-time.After(2 * time.Second):
			break
		}
	}, NoDelay(), Interval(3*time.Second))
	job.Start()

	// wait for start
	time.Sleep(time.Second)

	job.Stop()
	assert.True(t, canceled)
}

func Test_OnStopJobBetweenTicks_ShouldStopJobImmediately(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		done = true
	}, Delay(2*time.Second), Interval(3*time.Second))
	job.Start()

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	job.StopContext(ctx)
	cancel()
	assert.False(t, done)
}

func Test_OnJobStopDuringTick_ShouldWaitForJobGracefulShutdown(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		time.Sleep(2 * time.Second)
		done = true
	}, NoDelay(), Interval(3*time.Second))
	job.Start()

	// wait for start
	time.Sleep(time.Second)

	job.Stop()
	assert.True(t, done)
}

func Test_OnJobStopWithCancelledContext_ShouldStopWaitingForShutdown(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		time.Sleep(3 * time.Second)
		done = true
	}, NoDelay(), Interval(3*time.Second))
	job.Start()

	// wait for start
	time.Sleep(time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	job.StopContext(ctx)
	cancel()
	assert.False(t, done)
}

func Test_OnJobFinish_ShouldSendDoneSignal(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		done = true
	}, Delay(2*time.Second))
	job.Start()

	// wait for start
	time.Sleep(time.Second)

	<-job.Done()
	assert.True(t, done)
}
