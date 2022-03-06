package server

import (
	"github.com/hirpc/hrpc/configs"
	"github.com/hirpc/hrpc/life"
)

func (h HRPC) makeMessageQueue() error {
	for n, m := range h.opts.MQs {
		cfg, err := configs.Get().Get("messagequeues/" + n)
		if err != nil {
			return err
		}
		if err := m.Load(cfg); err != nil {
			return err
		}
		if err := m.Connect(); err != nil {
			return err
		}

		// in case of v has been overwriten
		mm := m
		life.WhenExit(mm.Destory)
	}
	return nil
}
