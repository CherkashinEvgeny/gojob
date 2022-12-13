package job

import (
	"time"
)

type Strategy interface {
	Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool)
}

func Compose(strategies ...Strategy) *CompositeStrategy {
	return &CompositeStrategy{
		strategies: strategies,
	}
}

type CompositeStrategy struct {
	strategies []Strategy
}

func (s *CompositeStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
	for index := 0; !ok && index < len(s.strategies); index++ {
		nextTickTime, ok = s.strategies[index].Tick(lastTickTime)
	}
	return
}

type StrategyFunc func(lastTickTime time.Time) (nextTickTime time.Time, ok bool)

func Function(f StrategyFunc) *FunctionStrategy {
	return &FunctionStrategy{
		f: f,
	}
}

type FunctionStrategy struct {
	f StrategyFunc
}

func (s *FunctionStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
	nextTickTime, ok = s.f(lastTickTime)
	return
}

func At(time time.Time) *ExactTimeStrategy {
	return &ExactTimeStrategy{
		time:      time,
		firstTick: true,
	}
}

type ExactTimeStrategy struct {
	time      time.Time
	firstTick bool
}

func (s *ExactTimeStrategy) Tick(_ time.Time) (nextTickTime time.Time, ok bool) {
	ok = s.firstTick
	if !ok {
		return
	}
	nextTickTime = s.time
	s.firstTick = false
	return
}

func NoDelay() *NoDelayStrategy {
	return &NoDelayStrategy{
		firstTick: true,
	}
}

type NoDelayStrategy struct {
	firstTick bool
}

func (s *NoDelayStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
	ok = s.firstTick
	if !ok {
		return
	}
	nextTickTime = lastTickTime
	s.firstTick = false
	return
}

func Delay(delay time.Duration) *DelayStrategy {
	return &DelayStrategy{
		delay:     delay,
		firstTick: true,
	}
}

type DelayStrategy struct {
	delay     time.Duration
	firstTick bool
}

func (s *DelayStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
	ok = s.firstTick
	if !ok {
		return
	}
	nextTickTime = lastTickTime.Add(s.delay)
	s.firstTick = false
	return
}

func Interval(period time.Duration) *IntervalStrategy {
	return IntervalIncludingPayloadDelay(period)
}

type IntervalStrategy = IntervalIncludingPayloadDelayStrategy

func IntervalIncludingPayloadDelay(period time.Duration) *IntervalIncludingPayloadDelayStrategy {
	return &IntervalIncludingPayloadDelayStrategy{
		period: period,
	}
}

type IntervalIncludingPayloadDelayStrategy struct {
	period time.Duration
}

func (s *IntervalIncludingPayloadDelayStrategy) Tick(_ time.Time) (nextTickTime time.Time, ok bool) {
	ok = true
	nextTickTime = time.Now().Add(s.period)
	return
}

func IntervalExcludingPayloadDelay(period time.Duration) *IntervalExcludingPayloadDelayStrategy {
	return &IntervalExcludingPayloadDelayStrategy{
		period: period,
	}
}

type IntervalExcludingPayloadDelayStrategy struct {
	period time.Duration
}

func (s *IntervalExcludingPayloadDelayStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
	ok = true
	nextTickTime = lastTickTime.Add(s.period)
	return
}
