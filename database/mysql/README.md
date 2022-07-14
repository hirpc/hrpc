# Introduction


## Usages

Customized use case.
```go
func main() {
    s, err := hrpc.NewServer(
		option.WithServerName("orderservice"),
		option.WithDatabases(mysql.New(
            mysql.WithCustomized(), // Very important, if you forget to provide this option, the HRPC will read configuration got from the configuration center to load.
            mysql.WithAddress("xxxx"),
            mysql.WithDB("aaa_db"),
            mysql.WithAuth("user", "pass"),
            mysql.WithPort(12144), // Default 3306
        )),
		option.WithEnvironment(option.Development),
		option.WithHealthCheck(),
	)
    // xxx
}
```


Customized use case with connection pool
```go
func main() {
    s, err := hrpc.NewServer(
		option.WithServerName("orderservice"),
		option.WithDatabases(mysql.New(
            mysql.WithCustomized(),
            mysql.WithAddress("xxxx"),
            mysql.WithDB("aaa_db"),
            mysql.WithAuth("user", "pass"),
            mysql.WithPort(12144), // Default 3306
            mysql.WithMaxOpenConns(10), // Default 3
            mysql.WithMaxIdleConns(4), // Default 1
        )),
		option.WithEnvironment(option.Development),
		option.WithHealthCheck(),
	)
    // xxx
}
```