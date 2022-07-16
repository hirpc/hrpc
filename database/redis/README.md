# Introduction


## Usages

Customized use case.
```go
func main() {
    s, err := hrpc.NewServer(
		option.WithServerName("orderservice"),
		option.WithDatabases(redis.New(
            redis.WithCustomized(), // Very important, if you forget to provide this option, the HRPC will read configuration got from the configuration center to load.
            redis.WithAddress("xxxx"),
            redis.WithDB(0),
            redis.WithAuth("user", "pass"),
            redis.WithPort(6378), // Default 6379
        )),
		option.WithEnvironment(option.Development),
		option.WithHealthCheck(),
	)
    // xxx
}
```
