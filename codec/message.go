package codec

import (
	"context"
	"encoding/json"
	"time"

	"github.com/hirpc/hrpc/utils/hash"
	"google.golang.org/grpc/metadata"
)

type MDKey string

func (m MDKey) String() string {
	return string(m)
}

const (
	MDMessage MDKey = "HRPC_MD_MESSAGE"
)

type MSG interface {
	Context() context.Context

	TraceID() string
	WithTraceID(s string)

	RequestTimeout() time.Duration
	WithRequestTimeout(t time.Duration)

	Namespace() string
	WithNamespace(n string)

	ServerName() string
	WithServerName(n string)

	Client() ClientInfo
	WithClient(ip, ua string)

	Metadata() metadata.MD
}

type message struct {
	context    context.Context `json:"-"`
	ServerN    string          `json:"server_name"`
	TID        string          `json:"trace_id"`
	RTimeout   time.Duration   `json:"request_timeout"`
	NameSpace  string          `json:"namespace"`
	ClientInfo ClientInfo      `json:"client_info"`
}

type ClientInfo struct {
	IP string `json:"ip"`
	UA string `json:"ua"`
}

func (m *message) Context() context.Context {
	return m.context
}

func (m *message) Client() ClientInfo {
	return m.ClientInfo
}

func (m *message) WithClient(ip, ua string) {
	m.ClientInfo.IP = ip
	m.ClientInfo.UA = ua
}

func (m *message) TraceID() string {
	return m.TID
}

func (m *message) RequestTimeout() time.Duration {
	return m.RTimeout
}

func (m *message) Namespace() string {
	return m.NameSpace
}

func (m *message) WithTraceID(s string) {
	m.TID = s
}

func (m *message) WithRequestTimeout(t time.Duration) {
	m.RTimeout = t
}

func (m *message) WithNamespace(n string) {
	m.NameSpace = n
}

func (m *message) ServerName() string {
	return m.ServerN
}

func (m *message) WithServerName(n string) {
	m.ServerN = n
}

func (m *message) Metadata() metadata.MD {
	d, err := json.Marshal(m)
	if err != nil {
		return metadata.Pairs(MDMessage.String(), "{}")
	}
	return metadata.Pairs(MDMessage.String(), string(d))
}

func Message(ctx context.Context) MSG {
	defaultMsg := message{
		context:  ctx,
		RTimeout: time.Second * 3,
		TID:      hash.SHA256(time.Now().UnixNano()),
	}
	md, exist := metadata.FromIncomingContext(ctx)
	if !exist {
		return &defaultMsg
	}
	vals := md.Get(MDMessage.String())
	if len(vals) == 0 {
		return &defaultMsg
	}
	var msg message
	if err := json.Unmarshal([]byte(vals[0]), &msg); err != nil {
		return &defaultMsg
	}
	msg.context = ctx
	return &msg
}
