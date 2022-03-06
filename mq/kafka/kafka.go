package kafka

import (
	"encoding/json"
	"log"
	"os"
	"strings"

	"github.com/Shopify/sarama"
	"github.com/hirpc/arsenal/uniqueid"
)

type Kafka struct {
	client sarama.Client
	opt    *Options
}

var k *Kafka

func New(opts ...Option) *Kafka {
	if k != nil {
		k.Destory()
	}
	k = &Kafka{
		opt: mergeOptions(opts...),
	}
	sarama.Logger = log.New(os.Stdout, "[sarama]", log.LstdFlags)
	return k
}

func (k *Kafka) Name() string {
	return "kafka"
}

// Load can load a set of configurations represented by a JSON string
// However, it is invalid if someone provides options at the contruction function
// Because the configurations provided at the New() function have the highest priority
func (k *Kafka) Load(src []byte) error {
	if k.opt.Brokers != "" {
		return nil
	}
	return json.Unmarshal(src, k.opt)
}

func (k *Kafka) Connect() error {
	config := sarama.NewConfig()
	config.ClientID = uniqueid.IP()
	config.Version = k.opt.version
	config.Producer.RequiredAcks = sarama.RequiredAcks(k.opt.Acks)
	config.Producer.Return.Successes = true
	config.Consumer.Offsets.AutoCommit.Enable = k.opt.AutoCommit
	config.Consumer.Offsets.Initial = k.opt.OffsetInitial
	if k.opt.SASL.enabled || k.opt.SASL.Username != "" {
		config.Net.SASL.Enable = true
		config.Net.SASL.User = k.opt.SASL.Username
		config.Net.SASL.Password = k.opt.SASL.Password
		config.Net.SASL.Mechanism = sarama.SASLTypePlaintext
	}

	c, err := sarama.NewClient(strings.Split(k.opt.Brokers, ","), config)
	if err != nil {
		return err
	}
	k.client = c
	return nil
}

func (k *Kafka) Destory() {
	if k.client != nil {
		k.client.Close()
	}
	return
}

func Client() sarama.Client {
	return k.client
}
