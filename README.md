# job

Golang library for job scheduling.

## About The Project

Job - small library, that provides api for job periodic execution.
Package has some benefits unlike known implementations:

- laconic api
- context package support
- easy customization

## Usage

Run with specified period:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Period(time.Second))
j.Start()
```

Run with specified interval (difference between period and interval is explained in [section](#underwater-rocks)):

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Interval(time.Second))
j.Start()
```

Run without initial delay:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Delay(0, job.Period(time.Second)))
j.Start()
```

Run with specified initial delay:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Delay(time.Second, job.Period(time.Second)))
j.Start()
```

Run at specific time:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.At(time.UnixMilli(1671040090920), job.Period(time.Second)))
j.Start()
```

Run yearly:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Yearly(time.February, 18, 10, 0, 0))
j.Start()
```

Run monthly:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Monthly(15, 10, 30, 0))
j.Start()
```

Run every last day of the month:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Monthly(-1, 0, 0, 0))
j.Start()
```

Run weekly:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Weekly(time.Monday, 11, 0, 0))
j.Start()
```

Run daily:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Daily(23, 30, 0))
j.Start()
```

Run hourly:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Hourly(30, 0))
j.Start()
```

Making complex timetable:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Timetable(
	job.Monthly(10, 10, 0, 0, 0),
	job.Monthly(25, 10, 0, 0, 0),
))
j.Start()
```

Using execution context:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Period(time.Second))
j.StartContext(ctx)
```

Stop job:

```
j := job.New(func(ctx context.Context) {
	fmt.Println("knock, knock (:")
}, job.Period(time.Second))
j.Start()
//...
j.Stop()
```

## Underwater rocks

### Interval vs Period

Interval - duration between last tick and next tick including payload execution time.
Period - duration between last tick and next tick excluding payload execution time.

This differance can be shown on the timeline:

```
|-----delay-----||---payload---||---delay---||-----payload-----|
|-----------interval-----------||-----------interval-----------|

|----delay-----||-------payload-------||----delay-----||---payload---|
|----period----|                       |----period----|
```

### Job restarting

There is no special api for job restart.
To get around this limitation, we can create new `Job` instance with appropriate parameters and start it.

## Similar projects

- [gocron](https://github.com/go-co-op/gocron)

## License

Job is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENCE.md)
for the full license text.

## Contact

- Email: `cherkashin.evgeny.viktorovich@gmail.com`
- Telegram: `@evgeny_cherkashin`