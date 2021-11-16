package hook

import "github.com/sirupsen/logrus"

// Hook 日志上传器
type Hook interface {
	// Establish for making connections
	Establish() error
	Fire(entry *logrus.Entry) error
	Levels() []logrus.Level
}
