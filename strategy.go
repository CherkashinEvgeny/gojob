package job

import (
	"sync"
	"time"
)

type Strategy interface {
	Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool)
}

var compositeStrategyPool = &sync.Pool{
	New: func() any {
		return &CompositeStrategy{make([]Strategy, 10)}
	},
}

func Compose(strategies ...Strategy) *CompositeStrategy {
	strategy := compositeStrategyPool.Get().(*CompositeStrategy)
	strategy.strategies = append(strategy.strategies[:0], strategies...)
	return strategy
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

type TickAtStrategy struct {
	time      time.Time
	firstTick bool
}

func TickAt(time time.Time) *TickAtStrategy {
	return &TickAtStrategy{
		time:      time,
		firstTick: true,
	}
}

func (s *TickAtStrategy) Tick(_ time.Time) (nextTickTime time.Time, ok bool) {
	ok = s.firstTick
	if !ok {
		return
	}
	nextTickTime = s.time
	s.firstTick = false
	return
}

type NoDelayStrategy struct {
	firstTick bool
}

func NoDelay() *NoDelayStrategy {
	return &NoDelayStrategy{
		firstTick: true,
	}
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

type DelayStrategy struct {
	delay     time.Duration
	firstTick bool
}

func Delay(delay time.Duration) *DelayStrategy {
	return &DelayStrategy{
		delay:     delay,
		firstTick: true,
	}
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

type IntervalStrategy = IntervalIncludingPayloadDelayStrategy

func Interval(period time.Duration) *IntervalStrategy {
	return IntervalIncludingPayloadDelay(period)
}

type IntervalIncludingPayloadDelayStrategy struct {
	period time.Duration
}

func IntervalIncludingPayloadDelay(period time.Duration) *IntervalIncludingPayloadDelayStrategy {
	return &IntervalIncludingPayloadDelayStrategy{
		period: period,
	}
}

func (s *IntervalIncludingPayloadDelayStrategy) Tick(_ time.Time) (nextTickTime time.Time, ok bool) {
	nextTickTime = time.Now().Add(s.period)
	ok = true
	return
}

type IntervalExcludingPayloadDelayStrategy struct {
	period time.Duration
}

func IntervalExcludingPayloadDelay(period time.Duration) *IntervalExcludingPayloadDelayStrategy {
	return &IntervalExcludingPayloadDelayStrategy{
		period: period,
	}
}

func (s *IntervalExcludingPayloadDelayStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time, ok bool) {
	nextTickTime = lastTickTime.Add(s.period)
	ok = true
	return
}
