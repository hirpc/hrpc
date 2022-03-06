package mq

type MQ interface {
	Load(src []byte) error
	Connect() error
	Name() string
	Destory()
}
