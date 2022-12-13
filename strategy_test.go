package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_OnCompositeStrategyTick_ShouldInquireUnderlyingStrategiesUntilSuccessAndReturnSuccessResult(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	counter := 0
	strategy := Compose(Function(func(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
		counter += 1
		ok = false
		return
	}), Function(func(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
		counter += 1
		ok = true
		nextTickTime = lastTickTime.Add(time.Second)
		return
	}), Function(func(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
		counter += 1
		ok = false
		return
	}))
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, 2, counter)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
	assert.True(t, ok)
}

func Test_OnFunctionStrategyTick_ShouldInvokeUnderlyingFunctionAndReturnItResult(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	counter := 0
	strategy := Function(func(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
		counter += 1
		ok = true
		nextTickTime = lastTickTime.Add(time.Second)
		return
	})
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, 1, counter)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
	assert.True(t, ok)
}

func Test_OnExactTimeStrategyFirstTick_ShouldReturnSpecifiedTime(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	tickTime := time.Now().Add(time.Second)
	strategy := At(tickTime)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, tickTime, nextTickTime)
	assert.True(t, ok)
}

func Test_OnExactTimeStrategySecondTick_ShouldReturnZeroTime(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	tickTime := time.Now().Add(time.Second)
	strategy := At(tickTime)
	_, _ = strategy.Tick(lastTickTime)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.True(t, nextTickTime.IsZero())
	assert.False(t, ok)
}

func Test_OnNoDelayStrategyFirstTick_ShouldReturnLastTickTime(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := NoDelay()
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime, nextTickTime)
	assert.True(t, ok)
}

func Test_OnNoDelayStrategySecondTick_ShouldReturnZeroTime(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := NoDelay()
	_, _ = strategy.Tick(lastTickTime)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.True(t, nextTickTime.IsZero())
	assert.False(t, ok)
}

func Test_OnDelayStrategyFirstTick_ShouldReturnLastTickTimeWithSpecifiedDelay(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Delay(time.Second)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
	assert.True(t, ok)
}

func Test_OnDelayStrategySecondTick_ShouldReturnZeroTime(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Delay(time.Second)
	_, _ = strategy.Tick(lastTickTime)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.True(t, nextTickTime.IsZero())
	assert.False(t, ok)
}

func Test_OnIntervalIncludingPayloadDelayStrategyFirstTick_ShouldReturnNowWithSpecifiedDelay(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := IntervalIncludingPayloadDelay(time.Second)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.True(t, time.Now().Add(time.Second).Sub(nextTickTime) > 0)
	assert.True(t, ok)
}

func Test_OnIntervalIncludingPayloadDelayStrategySecondTick_ShouldReturnNowWithSpecifiedDelay(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := IntervalIncludingPayloadDelay(time.Second)
	_, _ = strategy.Tick(lastTickTime)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.True(t, time.Now().Add(time.Second).Sub(nextTickTime) > 0)
	assert.True(t, ok)
}

func Test_OnIntervalExcludingPayloadDelayStrategyFirstTick_ShouldReturnLastTickTimeWithSpecifiedDelay(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := IntervalExcludingPayloadDelay(time.Second)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
	assert.True(t, ok)
}

func Test_OnIntervalExcludingPayloadDelayStrategySecondTick_ShouldReturnLastTickTimeWithSpecifiedDelay(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := IntervalExcludingPayloadDelay(time.Second)
	_, _ = strategy.Tick(lastTickTime)
	nextTickTime, ok := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
	assert.True(t, ok)
}
