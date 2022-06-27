package tracer

import (
	"context"
	"crypto/sha256"
	"fmt"
	"math/rand"
	"time"

	"github.com/hirpc/hrpc/codec"
	"github.com/hirpc/hrpc/log"
	"github.com/hirpc/hrpc/utils/hash"
	"github.com/hirpc/hrpc/utils/uniqueid"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
)

func prefix() string {
	seed := time.Now().UnixNano()
	rand.Seed(seed)
	h := sha256.New()
	h.Write([]byte(fmt.Sprintf("%v-%v", rand.Intn(99999), seed)))
	return fmt.Sprintf("%x", h.Sum(nil))
}

// NewID generates a random trace id in string
func NewID(serverName string) string {
	if serverName != "" {
		return hash.SHA256(serverName) + "." + uniqueid.String()
	}
	return prefix() + "." + uniqueid.String()
}

// AddTraceID will add an unique id to the ctx
func AddTraceID(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	msg := codec.Message(ctx)
	traceid := NewID(msg.ServerName())
	if msg.TraceID() == "" {
		msg.WithTraceID(traceid)
	}
	ctx, cancel := context.WithTimeout(msg.Context(), msg.RequestTimeout())
	defer cancel()

	v := msg.Metadata()
	ctx = metadata.NewOutgoingContext(
		ctx, v,
	)

	v1 := make(chan interface{}, 1)
	v2 := make(chan error, 1)
	v3 := make(chan time.Duration, 1)

	go func(ctx context.Context) {
		defer func() {
			if err := recover(); err != nil {
				log.WithFields(ctx).Error(err)
			}
		}()
		now := time.Now()
		resp, err := handler(ctx, req)
		v1 <- resp
		v2 <- err
		v3 <- time.Since(now)
	}(ctx)

	select {
	case <-ctx.Done():
		msg.WithRequestTimeout(
			time.Duration(0),
		)
		log.WithFields(ctx).Info("timeout")
		return nil, status.Error(codes.DeadlineExceeded, "deadline exceeded")
	case d := <-v3:
		left := time.Duration(
			(msg.RequestTimeout().Milliseconds() - d.Milliseconds()) * int64(time.Millisecond),
		)
		msg.WithRequestTimeout(left)
		log.WithFields(ctx, "used", d.String(), "left", left.String()).Info()
		return <-v1, <-v2
	}
}
