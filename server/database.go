package server

import (
	"errors"

	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/database/category"
	"github.com/hirpc/hrpc/life"
)

var ErrInvalidCategory = errors.New("invalid database category")

func (h HRPC) makeDatabase() error {
	for k, v := range h.opts.DBs {
		var cfg []byte
		var err error

		switch k {
		case category.MySQL:
			cfg, err = configs.Get().Get("databases/mysql")
			if err != nil {
				return err
			}
		case category.Redis:
			cfg, err = configs.Get().Get("databases/redis")
			if err != nil {
				return err
			}
		default:
			return ErrInvalidCategory
		}
		if err := v.Load(cfg); err != nil {
			return err
		}
		if err := v.Connect(); err != nil {
			return err
		}
		// in case of v has been overwriten
		vv := v
		life.WhenExit(vv.Destory)
	}
	return nil
}
