package server

import (
	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/life"
)

func (h HRPC) makeDatabase() error {
	for _, v := range h.opts.DBs {
		cfg, err := configs.Get().Get("databases/" + v.Name())
		if err != nil {
			return err
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
