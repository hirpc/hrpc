package category

type Category string

func (c Category) String() string {
	return string(c)
}

const (
	MySQL Category = "mysql"
	Redis Category = "redis"
)
