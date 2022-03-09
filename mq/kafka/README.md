# Introduction

## Usage

**For Consumer:**
```
import (
    // ...
    "github.com/hirpc/hrpc/mq/kafka"
    "github.com/hirpc/hrpc/log"
)
type tmp struct{}

func (t tmp) Handle(ctx context.Context, topic string, key, value []byte, partition int32, offset int64) error {
	log.WithFields(ctx, "KKK", "Kafka").Warn(topic, string(value), partition, offset)
	return nil
}


func main() {
    s, err := hrpc.NewServer(
        // ...
		option.WithMessageQueues(kafka.New(kafka.WithVersion("1.1.1"))),
        // ...
		option.WithHealthCheck(),
	)
    if err != nil {
        panic(err)
    }

    kafka.RegisterGroupConsumer(hrpc.BackgroundContext(), tmp{}, "GROUP_NAME", "TOPIC_NAME")

    if err := s.Serve(); err != nil {
		panic(err)
	}
}
```