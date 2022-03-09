package hook

import (
	"encoding/json"

	"github.com/hirpc/arsenal/plugins/tlog"
	"github.com/hirpc/hrpc/configs"
)

type cls struct {
	topic      string
	Credential tlog.Credential `json:"credential"`
	Endpoint   string          `json:"endpoint"`
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
	if err := json.Unmarshal([]byte(cfg), c); err != nil {
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
	if clslog.Endpoint != "" {
		return tlog.New(clslog.topic, clslog.Credential, tlog.WithEndpoint(clslog.Endpoint))
	}
	return tlog.New(clslog.topic, clslog.Credential)
}
