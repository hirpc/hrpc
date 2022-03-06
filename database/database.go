package database

type Database interface {
	Load(src []byte) error
	Connect() error
	Name() string
	Destory()
}
