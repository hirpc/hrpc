package database

import (
	"github.com/hirpc/hrpc/database/category"
)

type Database interface {
	Load(src []byte) error
	Connect() error
	Category() category.Category
	Destory()
}
