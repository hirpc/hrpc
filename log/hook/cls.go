package hook

import (
	"encoding/json"

	"github.com/hirpc/arsenal/plugins/tlog"
	"github.com/hirpc/hrpc/configs"
)

type cls struct {
	topic      string
	credential tlog.Credential
}

var clslog *cls

func NewCLS(topic string) *cls {
	clslog = &cls{
		topic: topic,
	}
	return clslog
}

func (c *cls) Load() error {
	cfg, err := configs.Get().Get("tencent/cls")
	if err != nil {
		return err
	}
	if err := json.Unmarshal([]byte(cfg), &c.credential); err != nil {
		return err
	}
	return nil
}

func (c *cls) Name() string {
	return "hrpc-cls"
}

func (c *cls) DependsOn() []string {
	return []string{"hrpc-configs"}
}

func CLSHook() *tlog.TLog {
	if clslog == nil {
		return nil
	}
	return tlog.New(clslog.topic, clslog.credential)
}
