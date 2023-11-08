package job

import (
	"time"
)

type Strategy interface {
	Tick(lastTickTime time.Time) (nextTickTime time.Time)
}

type StrategyFunc func(lastTickTime time.Time) (nextTickTime time.Time)

func Function(f StrategyFunc) FunctionStrategy {
	return FunctionStrategy{
		f: f,
	}
}

var _ Strategy = (*FunctionStrategy)(nil)

type FunctionStrategy struct {
	f StrategyFunc
}

func (s FunctionStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return s.f(lastTickTime)
}

func Delay(delay time.Duration, strategy Strategy) *DelayStrategy {
	return &DelayStrategy{
		applied:  false,
		delay:    delay,
		strategy: strategy,
	}
}

var _ Strategy = (*DelayStrategy)(nil)

type DelayStrategy struct {
	applied  bool
	delay    time.Duration
	strategy Strategy
}

func (s *DelayStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	if !s.applied {
		s.applied = true
		return time.Now().Add(s.delay)
	}
	return s.strategy.Tick(lastTickTime)
}

func At(time time.Time, strategy Strategy) *AtStrategy {
	return &AtStrategy{
		applied:  false,
		time:     time,
		strategy: strategy,
	}
}

var _ Strategy = (*AtStrategy)(nil)

type AtStrategy struct {
	applied  bool
	time     time.Time
	strategy Strategy
}

func (s *AtStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	if !s.applied {
		s.applied = true
		return s.time
	}
	return s.strategy.Tick(lastTickTime)
}

func Interval(interval time.Duration) IntervalStrategy {
	return IntervalStrategy{
		interval: interval,
	}
}

var _ Strategy = (*IntervalStrategy)(nil)

type IntervalStrategy struct {
	interval time.Duration
}

func (s IntervalStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return lastTickTime.Add(s.interval)
}

func Period(period time.Duration) PeriodStrategy {
	return PeriodStrategy{
		period: period,
	}
}

var _ Strategy = (*PeriodStrategy)(nil)

type PeriodStrategy struct {
	period time.Duration
}

func (s PeriodStrategy) Tick(_ time.Time) (nextTickTime time.Time) {
	return time.Now().Add(s.period)
}

func Timetable(timetable ...Strategy) TimetableStrategy {
	return TimetableStrategy{
		timetable: timetable,
	}
}

var _ Strategy = (*TimetableStrategy)(nil)

type TimetableStrategy struct {
	timetable []Strategy
}

func (s TimetableStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	for _, strategy := range s.timetable {
		strategyNextTickTime := strategy.Tick(lastTickTime)
		if nextTickTime.IsZero() || !strategyNextTickTime.IsZero() && nextTickTime.After(strategyNextTickTime) {
			nextTickTime = strategyNextTickTime
		}
	}
	return
}

func Yearly(month time.Month, day int, hour int, minute int, second int) YearlyStrategy {
	return YearlyStrategy{
		month:  month,
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	}
}

var _ Strategy = (*YearlyStrategy)(nil)

type YearlyStrategy struct {
	month  time.Month
	day    int
	hour   int
	minute int
	second int
}

func (s YearlyStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return nextYearPeriod(lastTickTime, s.month, s.day, s.hour, s.minute, s.second)
}

func Monthly(day int, hour int, minute int, second int) MonthlyStrategy {
	return MonthlyStrategy{
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	}
}

var _ Strategy = (*MonthlyStrategy)(nil)

type MonthlyStrategy struct {
	day    int
	hour   int
	minute int
	second int
}

func (s MonthlyStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return nextMonthPeriod(lastTickTime, s.day, s.hour, s.minute, s.second)
}

func Weekly(day time.Weekday, hour int, minute int, second int) WeeklyStrategy {
	return WeeklyStrategy{
		day:    day,
		hour:   hour,
		minute: minute,
		second: second,
	}
}

var _ Strategy = (*WeeklyStrategy)(nil)

type WeeklyStrategy struct {
	day    time.Weekday
	hour   int
	minute int
	second int
}

func (s WeeklyStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return nextWeekPeriod(lastTickTime, s.day, s.hour, s.minute, s.second)
}

func Daily(hour int, minute int, second int) DailyStrategy {
	return DailyStrategy{
		hour:   hour,
		minute: minute,
		second: second,
	}
}

var _ Strategy = (*DailyStrategy)(nil)

type DailyStrategy struct {
	hour   int
	minute int
	second int
}

func (s DailyStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return nextDayPeriod(lastTickTime, s.hour, s.minute, s.second)
}

func Hourly(minute int, second int) HourlyStrategy {
	return HourlyStrategy{
		minute: minute,
		second: second,
	}
}

var _ Strategy = (*HourlyStrategy)(nil)

type HourlyStrategy struct {
	minute int
	second int
}

func (s HourlyStrategy) Tick(lastTickTime time.Time) (nextTickTime time.Time) {
	return nextHourPeriod(lastTickTime, s.minute, s.second)
}

func nextYearPeriod(dt time.Time, month time.Month, day int, hour int, minute int, second int) (ndt time.Time) {
	year := dt.Year()
	location := dt.Location()
	for !ndt.After(dt) {
		ndt = time.Date(year, month, 1, hour, minute, second, 0, location)
		ndt = ndt.AddDate(0, 0, convertAbstractDayToDayNumber(ndt, day)-1)
		year++
	}
	return ndt
}

func nextMonthPeriod(dt time.Time, day int, hour int, minute int, second int) (ndt time.Time) {
	year := dt.Year()
	month := dt.Month()
	location := dt.Location()
	for !ndt.After(dt) {
		ndt = time.Date(year, month, 1, hour, minute, second, 0, location)
		ndt = ndt.AddDate(0, 0, convertAbstractDayToDayNumber(ndt, day)-1)
		month++
	}
	return ndt
}

func convertAbstractDayToDayNumber(dt time.Time, abstractDay int) (day int) {
	if abstractDay < 0 {
		return max(dayCountInCurrentMonth(dt)+abstractDay+1, 1)
	}
	if abstractDay == 0 {
		return 1
	}
	return min(abstractDay, dayCountInCurrentMonth(dt))
}

func max(a int, b int) (c int) {
	if a > b {
		return a
	}
	return b
}

func min(a int, b int) (c int) {
	if a < b {
		return a
	}
	return b
}

func dayCountInCurrentMonth(dt time.Time) (days int) {
	firstDayInMonth := time.Date(dt.Year(), dt.Month(), 1, 0, 0, 0, 0, dt.Location())
	lastDayInMonth := firstDayInMonth.AddDate(0, 1, -1)
	days = lastDayInMonth.Day()
	return days
}

func nextWeekPeriod(dt time.Time, weekday time.Weekday, hour int, minute int, second int) (ndt time.Time) {
	ndt = nextDayPeriod(dt, hour, minute, second)
	ndt = ndt.AddDate(0, 0, int(weekday-ndt.Weekday()))
	for !ndt.After(dt) {
		ndt = ndt.AddDate(0, 0, 7)
	}
	return ndt
}

func nextDayPeriod(dt time.Time, hour int, minute int, second int) (ndt time.Time) {
	ndt = nextHourPeriod(dt, minute, second)
	duration := time.Hour * time.Duration(hour-ndt.Hour())
	duration %= 24 * time.Hour
	if duration < 0 {
		duration += 24 * time.Hour
	}
	ndt = ndt.Add(duration)
	return ndt
}

func nextHourPeriod(dt time.Time, minute int, second int) (ndt time.Time) {
	ndt = nextMinutePeriod(dt, second)
	duration := time.Minute * time.Duration(minute-ndt.Minute())
	duration %= time.Hour
	if duration < 0 {
		duration += time.Hour
	}
	ndt = ndt.Add(duration)
	return ndt
}

func nextMinutePeriod(dt time.Time, second int) (ndt time.Time) {
	ndt = nextSecond(dt)
	duration := time.Second * time.Duration(second-ndt.Second())
	duration %= time.Minute
	if duration < 0 {
		duration += time.Minute
	}
	ndt = ndt.Add(duration)
	return ndt
}

func nextSecond(dt time.Time) (ndt time.Time) {
	duration := time.Second - (time.Duration(dt.Nanosecond()) % time.Second)
	ndt = dt.Add(duration)
	return ndt
}
