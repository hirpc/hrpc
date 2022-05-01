package log

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/hirpc/hrpc/codec"
	"github.com/hirpc/hrpc/utils/location"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func init() {
	if err := location.New().Load(); err != nil {
		panic(err)
	}
	setup()
}

var (
	logger = logrus.New()
	option = Option{
		Formatter:   &logrus.JSONFormatter{},
		Outer:       os.Stdout,
		Environment: "Unknown",
		StackSkip:   1,
	}
)

func setup() {
	logger.SetReportCaller(!option.DisableReportCaller)
	logger.SetFormatter(option.Formatter)
	for _, hook := range option.Hooks {
		if err := hook.Establish(); err != nil {
			fmt.Println("** establishing log failed, " + err.Error())
		}
		logger.AddHook(hook)
	}
	logger.Out = option.Outer
}

func WithFields(ctx context.Context, fields ...interface{}) *logrus.Entry {
	_, file, line, _ := runtime.Caller(option.StackSkip)
	msg := codec.Message(ctx)
	params := logrus.Fields{
		"category":    "system",
		"trace_id":    msg.TraceID(),
		"pst_time":    time.Now().In(location.PST()).Format("2006-01-02 15:04:05"),
		"bj_time":     time.Now().In(location.BJS()).Format("2006-01-02 15:04:05"),
		"file_name":   fmt.Sprintf("%s:%d", filepath.Base(file), line),
		"package":     filepath.Base(filepath.Dir(file)),
		"environment": option.Environment,
		"server_name": msg.ServerName(),
	}
	if len(fields)%2 != 0 {
		// fields should exist in pair
		return logger.WithFields(params)
	}
	for i := 0; i < len(fields); {
		key := "unknown"
		if k, ok := fields[i].(string); ok {
			key = k
		}
		params[key] = fields[i+1]
		i += 2
	}
	return logger.WithFields(params)
}

func AuditLog(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	msg := codec.Message(ctx)
	interceptorEntry(
		"trace_id", msg.TraceID(),
		"method", info.FullMethod,
		"direction", "incoming",
	).Infoln(req)
	resp, err := handler(ctx, req)
	interceptorEntry(
		"trace_id", msg.TraceID(),
		"method", info.FullMethod,
		"direction", "outcoming",
		"error", err,
	).Infoln(resp)
	return resp, err
}

func interceptorEntry(fields ...interface{}) *logrus.Entry {
	params := logrus.Fields{
		"category": "audit",
		"pst_time": time.Now().In(location.PST()).Format("2006-01-02 15:04:05"),
		"bj_time":  time.Now().In(location.BJS()).Format("2006-01-02 15:04:05"),
	}
	if len(fields)%2 != 0 {
		// fields should exist in pair
		return logger.WithFields(params)
	}
	for i := 0; i < len(fields); {
		key := "unknown"
		if k, ok := fields[i].(string); ok {
			key = k
		}
		params[key] = fields[i+1]
		i += 2
	}
	return logger.WithFields(params)
}
