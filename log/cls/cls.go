package cls

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/hirpc/hrpc/configs"
	"github.com/sirupsen/logrus"
	clssdk "github.com/tencentcloud/tencentcloud-cls-sdk-go"
)

type cls struct {
	Topic      string     `json:"topic_id"`
	Credential Credential `json:"credential"`
	Endpoint   string     `json:"endpoint"`

	producer *clssdk.AsyncProducerClient
	opt      Options
}

var (
	clslog *cls
	cb     *callback
)

func New(opts ...Option) *cls {
	opt := Options{
		// 默认入口
		endpoint: "na-siliconvalley.cls.tencentcs.com",
		callback: nil,
	}
	for _, o := range opts {
		o(&opt)
	}
	clslog = &cls{
		opt: opt,
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
	// option value has highest priority
	if c.opt.topicID != "" {
		c.Topic = c.opt.topicID
	}
	cb = &callback{
		c.opt.callback,
	}
	return nil
}

func (c *cls) Name() string {
	return "hrpc-cls"
}

func (c *cls) DependsOn() []string {
	return []string{"hrpc-configs"}
}

func (c *cls) Establish() error {
	producerConfig := clssdk.GetDefaultAsyncProducerClientConfig()
	producerConfig.Endpoint = c.opt.endpoint
	producerConfig.AccessKeyID = c.Credential.SecretID
	producerConfig.AccessKeySecret = c.Credential.SecretKey
	producerInstance, err := clssdk.NewAsyncProducerClient(producerConfig)
	if err != nil {
		fmt.Println("NewAsyncProducerClient failed, " + err.Error())
		return err
	}
	producerInstance.Start()
	c.producer = producerInstance
	return nil
}

func (c cls) Fire(entry *logrus.Entry) error {
	if err := c.producer.SendLog(
		c.Topic,
		clssdk.NewCLSLog(
			time.Now().Unix(),
			logContent(entry),
		), cb,
	); err != nil {
		fmt.Println("failed to send log, " + err.Error())
	}
	return nil
}

func logContent(entry *logrus.Entry) map[string]string {
	var out = map[string]string{
		"Level":   entry.Level.String(),
		"Message": entry.Message,
		"Time":    entry.Time.Format("2006-01-02 15:04:05"),
	}

	d, err := json.Marshal(entry.Data)
	if err == nil {
		out["Fields"] = string(d)
	}

	v, err := json.Marshal(entry.Caller)
	if err == nil {
		out["Caller"] = string(v)
	}
	return out
}

func (c cls) Levels() []logrus.Level {
	return []logrus.Level{
		logrus.PanicLevel,
		logrus.FatalLevel,
		logrus.ErrorLevel,
		logrus.WarnLevel,
		logrus.InfoLevel,
	}
}

func Hook() *cls {
	return clslog
}
