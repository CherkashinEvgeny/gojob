package job

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func Test_OnFunctionStrategyTick_ShouldInvokeUnderlyingFunctionAndReturnItResult(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(time.Second)
	})
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
}

func Test_OnDelayStrategyFirstTick_ShouldReturnTickTimeUsingSpecifiedDelay(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Delay(time.Second, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(time.Minute)
	}))
	nextTickTime := strategy.Tick(lastTickTime)
	assert.InDelta(t, time.Now().Add(time.Second).UnixNano(), nextTickTime.UnixNano(), float64(10*time.Millisecond))
}

func Test_OnDelayStrategySecondTick_ShouldReturnTickTimeUsingSpecifiedStrategy(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Delay(time.Second, Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(time.Minute)
	}))
	_ = strategy.Tick(lastTickTime)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Minute), nextTickTime)
}

func Test_OnAtStrategyFirstTick_ShouldReturnSpecifiedTickTime(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := At(lastTickTime.Add(time.Second), Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(time.Minute)
	}))
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
}

func Test_OnAtStrategySecondTick_ShouldReturnTickTimeUsingSpecifiedStrategy(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := At(lastTickTime.Add(time.Second), Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
		return lastTickTime.Add(time.Minute)
	}))
	_ = strategy.Tick(lastTickTime)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Minute), nextTickTime)
}

func Test_OnIntervalStrategyTick_ShouldReturnTickTimeAccordingToSpecifiedPeriod(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Interval(2 * time.Second)
	_ = strategy.Tick(lastTickTime)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(2*time.Second), nextTickTime)
}

func Test_OnPeriodStrategyTick_ShouldReturnTickTimeAccordingToSpecifiedPeriod(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Period(2 * time.Second)
	_ = strategy.Tick(lastTickTime)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.InDelta(t, time.Now().Add(2*time.Second).UnixNano(), nextTickTime.UnixNano(), float64(10*time.Millisecond))
}

func Test_OnTimetableStrategyTick_ShouldReturnSmallestTickTimeOfUnderlyingStrategies(t *testing.T) {
	lastTickTime := time.Now().Add(-2 * time.Second)
	strategy := Timetable(
		Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
			return lastTickTime.Add(time.Second)
		}),
		Function(func(lastTickTime time.Time) (nextTickTime time.Time) {
			return lastTickTime.Add(2 * time.Second)
		}),
	)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, lastTickTime.Add(time.Second), nextTickTime)
}

func Test_OnYearlyStrategyTickWhenLastTickTimeLessThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 11, 39, 2, 0, time.Local)
	strategy := Yearly(time.February, 18, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnYearlyStrategyTickWhenLastTickEqualTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local)
	strategy := Yearly(time.February, 18, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2024, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnYearlyStrategyTickWhenLastTickTimeGraterThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 19, 11, 39, 2, 0, time.Local)
	strategy := Yearly(time.February, 18, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2024, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnMonthlyStrategyTickWhenLastTickTimeLessThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 11, 39, 2, 0, time.Local)
	strategy := Monthly(18, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnMonthlyStrategyTickWhenLastTickEqualTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local)
	strategy := Monthly(18, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.March, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnMonthlyStrategyTickWhenLastTickTimeGraterThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 19, 11, 39, 2, 0, time.Local)
	strategy := Monthly(18, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.March, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnMonthlyStrategyTickWithNegativeDayNumber_ShouldReturnTickTimeOfTheReverseMonthlyStrategy(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 19, 11, 39, 2, 0, time.Local)
	strategy := Monthly(-2, 10, 0, 0)
	reverseStrategy := Monthly(27, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	reverseNextTickTime := reverseStrategy.Tick(lastTickTime)
	assert.Equal(t, reverseNextTickTime, nextTickTime)
}

func Test_OnMonthlyStrategyTickWithOverflowDayNumber_ShouldRoundDayToTheBiggestDayOfTheMonthAndReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 19, 11, 39, 2, 0, time.Local)
	strategy := Monthly(31, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 28, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnWeeklyStrategyTickWhenLastTickTimeLessThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 11, 39, 2, 0, time.Local)
	strategy := Weekly(time.Saturday, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnWeeklyStrategyTickWhenLastTickEqualTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local)
	strategy := Weekly(time.Saturday, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 25, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnWeeklyStrategyTickWhenLastTickTimeGraterThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 19, 11, 39, 2, 0, time.Local)
	strategy := Weekly(time.Saturday, 10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 25, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnDailyStrategyTickWhenLastTickTimeLessThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 9, 39, 2, 0, time.Local)
	strategy := Daily(10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 17, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnDailyStrategyTickWhenLastTickEqualTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 10, 0, 0, 0, time.Local)
	strategy := Daily(10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnDailyStrategyTickWhenLastTickTimeGraterThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 11, 39, 2, 0, time.Local)
	strategy := Daily(10, 0, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 18, 10, 0, 0, 0, time.Local), nextTickTime)
}

func Test_OnHourlyStrategyTickWhenLastTickTimeLessThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 9, 28, 2, 0, time.Local)
	strategy := Hourly(30, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 17, 9, 30, 0, 0, time.Local), nextTickTime)
}

func Test_OnHourlyStrategyTickWhenLastTickEqualTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 9, 30, 0, 0, time.Local)
	strategy := Hourly(30, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 17, 10, 30, 0, 0, time.Local), nextTickTime)
}

func Test_OnHourlyStrategyTickWhenLastTickTimeGraterThanTimeInConfiguration_ShouldReturnNearestTickTimeInTheFutureAccordingToConfiguration(t *testing.T) {
	lastTickTime := time.Date(2023, time.February, 17, 9, 39, 2, 0, time.Local)
	strategy := Hourly(30, 0)
	nextTickTime := strategy.Tick(lastTickTime)
	assert.Equal(t, time.Date(2023, time.February, 17, 10, 30, 0, 0, time.Local), nextTickTime)
}
