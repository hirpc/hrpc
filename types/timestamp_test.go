package types

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

func TestTimestamp_MarshalBSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  time.Time
		want    []byte
		wantErr bool
	}{
		{
			name:    "Test1",
			fields:  time.Unix(1643549328, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := temp{
				Create: Timestamp{tt.fields},
				Name:   "test1",
			}
			got, err := bson.Marshal(tr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			model := struct {
				Create time.Time
				Name   string
			}{
				Create: time.Unix(1643549328, 0),
				Name:   "test1",
			}
			tt.want, err = bson.Marshal(&model)
			if (err != nil) != tt.wantErr {
				t.Errorf("Time.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Timestamp.MarshalJSON() = %v, want %v", string(got), tt.want)
			}
		})
	}
}

func TestTimestamp_UnmarshalBSON(t *testing.T) {
	tests := []struct {
		name    string
		fields  time.Time
		want    time.Time
		wantErr bool
	}{
		{
			name:    "Test1",
			fields:  time.Unix(1643549328, 0),
			want:    time.Unix(1643549328, 0),
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tr := temp{
				Create: Timestamp{tt.fields},
				Name:   "test1",
			}
			marshalled, err := bson.Marshal(tr)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.MarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			var unmarshalled temp
			err = bson.Unmarshal(marshalled, &unmarshalled)
			if (err != nil) != tt.wantErr {
				t.Errorf("Timestamp.UnmarshalJSON() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !unmarshalled.Create.Time.Equal(tt.want) {
				t.Errorf("Timestamp.MarshalJSON() = %v, want %v", unmarshalled.Create.Time, tt.want)
			}
		})
	}
}

func TestTimestamp_Scan(t *testing.T) {
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
