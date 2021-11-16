package hook

import (
	"encoding/json"
	"runtime"
	"time"

	"github.com/Shopify/sarama"
	"github.com/sirupsen/logrus"
)

type kafka struct {
	user, password string
	endpoint       string
	topic          string

	producer sarama.SyncProducer
}

func NewKafka(user, password, endpoint, topic string) *kafka {
	return &kafka{
		user:     user,
		password: password,
		endpoint: endpoint,
		topic:    topic,
	}
}

func (k kafka) Fire(entry *logrus.Entry) error {
	if _, _, err := k.producer.SendMessage(&sarama.ProducerMessage{
		Topic: k.topic,
		Value: sarama.StringEncoder(getContent(entry)),
	}); err != nil {
		return err
	}
	return nil
}

func (k kafka) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

func (k *kafka) Establish() error {
	producer, err := sarama.NewSyncProducer([]string{k.endpoint}, k.config())
	if err != nil {
		return err
	}
	k.producer = producer
	return nil
}

func (k kafka) config() *sarama.Config {
	config := sarama.NewConfig()
	config.Net.SASL.Mechanism = "PLAIN"
	config.Net.SASL.Version = int16(1)
	config.Net.SASL.Enable = true
	config.Net.SASL.User = k.user
	config.Net.SASL.Password = k.password
	config.Producer.Return.Successes = true
	config.Version = sarama.V0_11_0_0
	return config
}

func getContent(entry *logrus.Entry) string {
	var out = struct {
		Fields  logrus.Fields
		Level   string
		Message string
		Time    time.Time
		Caller  *runtime.Frame
	}{
		Fields:  entry.Data,
		Level:   entry.Level.String(),
		Message: entry.Message,
		Time:    entry.Time,
		Caller:  entry.Caller,
	}
	d, err := json.Marshal(out)
	if err != nil {
		return ""
	}
	return string(d)
}
