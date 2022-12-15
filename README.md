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
}, job.NoDelay(), job.Period(time.Second))
j.Start()
```

Run with specified initial delay:

```
j := job.New(func(ctx context.Context) {
    fmt.Println("knock, knock (:")
}, job.Delay(time.Second), job.Period(time.Second))
j.Start()
```

Run at specific time:

```
j := job.New(func(ctx context.Context) {
    fmt.Println("knock, knock (:")
}, job.At(time.UnixMilli(1671040090920)), job.Period(time.Second))
j.Start()
```

Using execution context:

```
j := job.New(func(ctx context.Context) {
    fmt.Println("knock, knock (:")
}, job.NoDelay(), job.Period(time.Second))
j.StartContext(ctx)
```

Stop job:

```
j := job.New(func(ctx context.Context) {
    fmt.Println("knock, knock (:")
}, job.At(time.UnixMilli(1671040090920)), job.Period(time.Second))
j.Start()
//...
j.Stop()
```

## Underwater rocks

### Arguments order

Be careful:

```
job.New(someFunc, job.NoDelay(), job.Period(time.Second))
```

not similar to

```
job.New(someFunc, job.Period(time.Second), job.NoDelay())
```

First block of code creates new job without initial delay and with one-second execution period.
In second code block, `job.NoDelay()` argument will be useless, because `job.Period(time.Second)` is infinite strategy.

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