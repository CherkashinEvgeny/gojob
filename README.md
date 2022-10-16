# job
Golang library for job scheduling.

## About The Project
Job - small library, that provides api for task scheduling.
Package has some benefits unlike known implementations:
- laconic api
- context package support

## Usage
Start job:
```
job := New(func(_ context.Context) {
	// your logic here
})
err := job.Start(0, time.Second)
```
Start job in context:
```
job := New(func(_ context.Context) {
	// your logic here
})
ctx, cancel := context.WithCancel(context.Background())
go func() {
    // wait some event
    cancel()
}()
err := job.StartContext(ctx, 0, time.Second)
```
Stop job:
```
err := job.Stop()
```
Stop job with context:
```
ctx, cancel := context.WithCancel(context.Background())
go func() {
    // wait some event
    cancel()
}()
err := job.StopContext(ctx)
```

## Similar projects
- [gocron](https://github.com/go-co-op/gocron)

## License
Retry is licensed under the Apache License, Version 2.0. See [LICENSE](./LICENCE.md) 
for the full license text.

## Contact
- Email: `cherkashin.evgeny.viktorovich@gmail.com`
- Telegram: `@evgeny_cherkashin`