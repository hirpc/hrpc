package kafka

import "github.com/Shopify/sarama"

type Options struct {
	Brokers       string              `json:"brokers"`
	SASL          sasl                `json:"sasl"`
	Acks          RequiredAcks        `json:"acks"`
	AutoCommit    bool                `json:"auto_commit"`
	OffsetInitial int64               `json:"offset_initial"`
	version       sarama.KafkaVersion `json:"-"`
}

type sasl struct {
	enabled  bool   `json:"-"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Option func(*Options)

func mergeOptions(opts ...Option) *Options {
	var opt = &Options{
		OffsetInitial: OffsetNewest,
		Acks:          WaitForLocal,
		version:       sarama.V1_1_1_0,
	}
	for _, o := range opts {
		o(opt)
	}
	return opt
}

func WithSASL(username, password string) Option {
	return func(o *Options) {
		o.SASL.enabled = true
		o.SASL.Username = username
		o.SASL.Password = password
	}
}

func WithAcks(ack RequiredAcks) Option {
	return func(o *Options) {
		o.Acks = ack
	}
}

func WithAutoCommit() Option {
	return func(o *Options) {
		o.AutoCommit = true
	}
}

func WithOffsetInitial(i int64) Option {
	return func(o *Options) {
		o.OffsetInitial = i
	}
}

func WithVersion(v string) Option {
	var version sarama.KafkaVersion
	switch v {
	case "1.1.1":
		version = sarama.V1_1_1_0
	default:
		version = sarama.V1_1_1_0
	}
	return func(o *Options) {
		o.version = version
	}
}

type RequiredAcks int16

const (
	// NoResponse doesn't send any response, the TCP ACK is all you get.
	NoResponse RequiredAcks = 0
	// WaitForLocal waits for only the local commit to succeed before responding.
	WaitForLocal RequiredAcks = 1
	// WaitForAll waits for all in-sync replicas to commit before responding.
	// The minimum number of in-sync replicas is configured on the broker via
	// the `min.insync.replicas` configuration key.
	WaitForAll RequiredAcks = -1
)

const (
	// OffsetNewest stands for the log head offset, i.e. the offset that will be
	// assigned to the next message that will be produced to the partition. You
	// can send this to a client's GetOffset method to get this offset, or when
	// calling ConsumePartition to start consuming new messages.
	OffsetNewest int64 = -1
	// OffsetOldest stands for the oldest offset available on the broker for a
	// partition. You can send this to a client's GetOffset method to get this
	// offset, or when calling ConsumePartition to start consuming from the
	// oldest offset that is still available on the broker.
	OffsetOldest int64 = -2
)
