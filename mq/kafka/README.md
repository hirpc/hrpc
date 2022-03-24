# Introduction

## Usage

**For Initialization:**
```
import (
    // ...
    "github.com/hirpc/hrpc/mq/kafka"
    "github.com/hirpc/hrpc/log"
    "github.com/hirpc/hrpc/option"
)

func main() {
    s, err := hrpc.NewServer(
		option.WithMessageQueues(kafka.New(kafka.WithVersion("1.1.1"))),
	)
    if err != nil {
        panic(err)
    }
    // ....
}
```

**For Producer:**
```
import (
    "github.com/hirpc/hrpc/log"
    "github.com/hirpc/hrpc/mq/kafka"
)

func Foo() error {
    if err := kafka.Produce(
		"TOPIC_NAME",
		*kafka.NewProduceMessage([]byte("MESSAGE")),
	); err != nil {
		log.WithFields(ctx).Error(err)
		return err
	}
	return nil
}

```

**For Consumer:**
```
import (
    // ...
    "github.com/hirpc/hrpc/mq/kafka"
    "github.com/hirpc/hrpc/log"
    "github.com/hirpc/hrpc/option"
)

type tmp struct{}

// Handle will receive messages from the kafka with topic, key, value, partition, offset/
// Every message received should return nil or error. The difference between these two values is:
//  nil   -> do commit for this message
//  error -> do NOT commit for this message
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