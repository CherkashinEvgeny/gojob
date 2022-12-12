package job

import (
	"context"
	"gotest.tools/assert"
	"testing"
	"time"
)

func Test_OnStartJobWithoutDelay_ShouldInvokePayloadFunctionImmediately(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		done = true
	})
	err := job.Start(NoDelay(), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	err = job.Stop()
	assert.NilError(t, err)
	assert.Check(t, done)
}

func Test_OnStartJobWithDelay_ShouldInvokePayloadFunctionWithDelay(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		done = true
	})
	err := job.Start(Delay(time.Second), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	//wait for delay
	time.Sleep(time.Second)
	assert.Check(t, done)

	err = job.Stop()
	assert.NilError(t, err)
}

func Test_OnStartJob_ShouldInvokePayloadFunctionPeriodically(t *testing.T) {
	counter := 0
	job := New(func(_ context.Context) {
		counter++
	})
	err := job.Start(NoDelay(), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	for i := 1; i < 6; i++ {
		assert.Equal(t, i, counter)
		time.Sleep(3 * time.Second)
	}
	err = job.Stop()
	assert.NilError(t, err)
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
	})
	err := job.Start(NoDelay(), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	err = job.Stop()
	assert.NilError(t, err)
	assert.Check(t, canceled)
}

func Test_OnStopJobWithDelay_ShouldStopImmediately(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		time.Sleep(2 * time.Second)
		done = true
	})
	err := job.Start(Delay(2*time.Second), Interval(3*time.Second))
	assert.NilError(t, err)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err = job.StopContext(ctx)
	cancel()
	assert.Check(t, !done)
	assert.NilError(t, err)
}

func Test_OnJobStop_ShouldWaitForJobGracefulShutdown(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		time.Sleep(2 * time.Second)
		done = true
	})
	err := job.Start(NoDelay(), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	err = job.Stop()
	assert.NilError(t, err)
	assert.Check(t, done)
}

func Test_OnJobStopAndContextTimeout_ShouldStopWaitingForShutdownAndReturnContextCancellationError(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		time.Sleep(3 * time.Second)
		done = true
	})
	err := job.Start(NoDelay(), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err = job.StopContext(ctx)
	cancel()
	assert.Check(t, !done)
	assert.Equal(t, ctx.Err(), err)
}

func Test_OnJobStop_ShouldSendDoneSignalTo(t *testing.T) {
	done := false
	job := New(func(_ context.Context) {
		time.Sleep(3 * time.Second)
		done = true
	})
	err := job.Start(NoDelay(), Interval(3*time.Second))
	assert.NilError(t, err)

	// wait for start
	time.Sleep(time.Second)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	err = job.StopContext(ctx)
	cancel()
	assert.Check(t, !done)
	assert.Equal(t, ctx.Err(), err)
}
