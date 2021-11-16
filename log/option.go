package log

import (
	"io"

	"github.com/hirpc/hrpc/log/hook"
	"github.com/sirupsen/logrus"
)

// Option defines a set of options
type Option struct {
	// DisableReportCaller will remove the report caller info if true
	DisableReportCaller bool
	// Formatter will format the log message as output
	Formatter logrus.Formatter
	// Outer defines where the log info outputs
	Outer io.Writer
	// Hooks 钩子回调
	Hooks []hook.Hook
	// Environment the env string
	Environment string
	StackSkip   int
}

// With should be used if you want to customize you logger
func With(opts ...Option) {
	for _, opt := range opts {
		if opt.DisableReportCaller {
			option.DisableReportCaller = opt.DisableReportCaller
		}
		if opt.Formatter != nil {
			option.Formatter = opt.Formatter
		}
		if opt.Outer != nil {
			option.Outer = opt.Outer
		}
		if len(opt.Hooks) != 0 {
			option.Hooks = append(option.Hooks, opt.Hooks...)
		}
		if opt.Environment != "" {
			option.Environment = opt.Environment
		}
		if opt.StackSkip != 0 {
			option.StackSkip = opt.StackSkip
		}
	}
	// update settings
	setup()
}

// AddHook for adding a hook
func AddHooks(hooks ...hook.Hook) error {
	if logger == nil {
		return nil
	}
	for _, hook := range hooks {
		if err := hook.Establish(); err != nil {
			return err
		}
		logger.AddHook(hook)
	}
	return nil
}
