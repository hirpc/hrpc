package types

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"
)

type temp struct {
	Create Timestamp
	Name   string
}

func TestTimestamp_MarshalJSON(t *testing.T) {
	type fields struct {
		Time time.Time
	}
	tests := []struct {
		name    string
		fields  fields
		want    []byte
		wantErr bool
	}{
		{
			name:    "Test1",
			fields:  fields{Time: time.Unix(1643549328, 0)},
			want:    []byte(`{"Create":1643549328,"Name":"test1"}`),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := temp{
				Create: Timestamp{tt.fields.Time},
				Name:   "test1",
			}

			got, err := json.Marshal(tr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Timestamp.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestTimestamp_Scan(t *testing.T) {
	type fields struct {
		Time time.Time
	}
	type args struct {
		src interface{}
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Test1",
			args:    args{src: []byte("1643549035")},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := Timestamp{}
			if err := tr.Scan(tt.args.src); (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.Scan() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
