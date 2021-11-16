package codec

import (
	"context"
	"reflect"
	"testing"
	"time"

	"google.golang.org/grpc/metadata"
)

func Test_message_Metadata(t *testing.T) {
	type fields struct {
		context   context.Context
		TID       string
		RTimeout  time.Duration
		NameSpace string
	}
	tests := []struct {
		name   string
		fields fields
		want   metadata.MD
	}{
		{
			name: "Test1",
			fields: fields{
				context:   context.Background(),
				TID:       "abcd",
				RTimeout:  time.Second,
				NameSpace: "Dev",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			m := &message{
				context:   tt.fields.context,
				TID:       tt.fields.TID,
				RTimeout:  tt.fields.RTimeout,
				NameSpace: tt.fields.NameSpace,
			}
			got := m.Metadata()
			t.Log(got.Get(MDMessage.String()))
			if got := m.Metadata(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("message.Metadata() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMessage(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name string
		args args
		want MSG
	}{
		{
			name: "Test1",
			args: args{
				ctx: metadata.NewIncomingContext(
					context.Background(), metadata.Pairs(MDMessage.String(), `{"trace_id":"abcd","request_timeout":1000000000,"namespace":"Dev"}`),
				),
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := Message(tt.args.ctx)
			t.Log(got.TraceID())
			if got := Message(tt.args.ctx); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Message() = %v, want %v", got, tt.want)
			}
		})
	}
}
