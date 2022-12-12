package job

import (
	"context"
	"sync"
	"time"
)

type Strategy interface {
	Tick(ctx context.Context) (ok bool)
}

type Resetter interface {
	Reset()
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

func (s *CompositeStrategy) Tick(ctx context.Context) (ok bool) {
	for index := 0; !ok && index < len(s.strategies); index++ {
		ok = s.strategies[index].Tick(ctx)
	}
	return
}

func (s *CompositeStrategy) Reset() {
	for _, strategy := range s.strategies {
		if resetter, ok := strategy.(Resetter); ok {
			resetter.Reset()
		}
	}
	compositeStrategyPool.Put(s)
}

var startAtStrategyPool = &sync.Pool{
	New: func() any {
		return &StartAtStrategy{}
	},
}

type StartAtStrategy struct {
	time time.Time
}

func StartAt(time time.Time) *StartAtStrategy {
	strategy := startAtStrategyPool.Get().(*StartAtStrategy)
	strategy.time = time
	return strategy
}

func (s *StartAtStrategy) Tick(ctx context.Context) (ok bool) {
	delay := time.Until(s.time)
	if delay <= 0 {
		ok = true
		return
	}
	ok = Sleep(ctx, delay)
	return
}

func (s *StartAtStrategy) Reset() {
	startAtStrategyPool.Put(s)
	return
}

var noDelayStrategyPool = &sync.Pool{
	New: func() any {
		return &NoDelayStrategy{true}
	},
}

type NoDelayStrategy struct {
	firstTick bool
}

func NoDelay() *NoDelayStrategy {
	strategy := noDelayStrategyPool.Get().(*NoDelayStrategy)
	strategy.firstTick = true
	return strategy
}

func (s *NoDelayStrategy) Tick(_ context.Context) (ok bool) {
	ok = s.firstTick
	s.firstTick = false
	return
}

func (s *NoDelayStrategy) Reset() {
	noDelayStrategyPool.Put(s)
	return
}

var delayStrategyPool = &sync.Pool{
	New: func() any {
		return &DelayStrategy{true, 0}
	},
}

type DelayStrategy struct {
	firstTick bool
	delay     time.Duration
}

func Delay(delay time.Duration) *DelayStrategy {
	strategy := delayStrategyPool.Get().(*DelayStrategy)
	strategy.firstTick = true
	strategy.delay = delay
	return strategy
}

func (s *DelayStrategy) Tick(ctx context.Context) (ok bool) {
	if !s.firstTick {
		return false
	}
	s.firstTick = true
	ok = Sleep(ctx, s.delay)
	return
}

func (s *DelayStrategy) Reset() {
	delayStrategyPool.Put(s)
	return
}

var Interval = IntervalIncludingPayloadDelay

var intervalIncludingPayloadDelayStrategyPool = &sync.Pool{
	New: func() any {
		return &IntervalIncludingPayloadDelayStrategy{true, 0, time.Time{}, nil}
	},
}

type IntervalIncludingPayloadDelayStrategy struct {
	firstTick    bool
	period       time.Duration
	nextTickTime time.Time
	timer        *time.Timer
}

func IntervalIncludingPayloadDelay(period time.Duration) *IntervalIncludingPayloadDelayStrategy {
	strategy := intervalIncludingPayloadDelayStrategyPool.Get().(*IntervalIncludingPayloadDelayStrategy)
	strategy.firstTick = true
	strategy.period = period
	return strategy
}

func (s *IntervalIncludingPayloadDelayStrategy) Tick(ctx context.Context) (ok bool) {
	if s.firstTick {
		s.nextTickTime = time.Now().Add(s.period)
		s.timer = timerFromPool(time.Until(s.nextTickTime))
		s.firstTick = false
	}
	select {
	case <-ctx.Done():
		ok = false
	case <-s.timer.C:
		s.nextTickTime = s.nextTickTime.Add(s.period)
		s.timer.Reset(time.Until(s.nextTickTime))
		ok = true
	}
	return
}

func (s *IntervalIncludingPayloadDelayStrategy) Reset() {
	if !s.firstTick {
		returnTimerToPool(s.timer)
		s.timer = nil
	}
	intervalIncludingPayloadDelayStrategyPool.Put(s)
	return
}

var intervalExcludingPayloadDelayStrategyPool = &sync.Pool{
	New: func() any {
		return &IntervalExcludingPayloadDelayStrategy{true, 0, time.Time{}, nil}
	},
}

type IntervalExcludingPayloadDelayStrategy struct {
	firstTick    bool
	period       time.Duration
	nextTickTime time.Time
	timer        *time.Timer
}

func IntervalExcludingPayloadDelay(period time.Duration) *IntervalExcludingPayloadDelayStrategy {
	strategy := intervalExcludingPayloadDelayStrategyPool.Get().(*IntervalExcludingPayloadDelayStrategy)
	strategy.firstTick = true
	strategy.period = period
	return strategy
}

func (s *IntervalExcludingPayloadDelayStrategy) Tick(ctx context.Context) (ok bool) {
	if s.firstTick {
		s.nextTickTime = time.Now().Add(s.period)
		s.timer = timerFromPool(time.Until(s.nextTickTime))
		s.firstTick = false
	} else {
		s.nextTickTime = s.nextTickTime.Add(s.period)
		s.timer.Reset(time.Until(s.nextTickTime))
	}
	select {
	case <-ctx.Done():
		ok = false
	case <-s.timer.C:
		ok = true
	}
	return
}

func (s *IntervalExcludingPayloadDelayStrategy) Reset() {
	if !s.firstTick {
		returnTimerToPool(s.timer)
		s.timer = nil
	}
	intervalExcludingPayloadDelayStrategyPool.Put(s)
	return
}

func Sleep(ctx context.Context, delay time.Duration) (ok bool) {
	timer := timerFromPool(delay)
	defer returnTimerToPool(timer)
	select {
	case <-timer.C:
		ok = true
	case <-ctx.Done():
		ok = false
	}
	return
}

var timerPool = &sync.Pool{}

func timerFromPool(delay time.Duration) (timer *time.Timer) {
	timerFromPool := timerPool.Get()
	if timerFromPool == nil {
		timer = time.NewTimer(delay)
	} else {
		timer = timerFromPool.(*time.Timer)
		timer.Reset(delay)
	}
	return
}

func returnTimerToPool(timer *time.Timer) {
	if !timer.Stop() {
		select {
		case <-timer.C:
			break
		default:
			break
		}
	}
	timerPool.Put(timer)
	return
}
