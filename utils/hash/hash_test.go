package hash

import (
	"testing"
)

func TestSHA1(t *testing.T) {
	type args struct {
		d []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test1",
			args: args{
				d: []string{"aaa", "bbb"},
			},
			want: "68d8572c2662b0f06f723d7d507954fb038b8558",
		},
		{
			name: "Test2",
			args: args{
				d: []string{},
			},
			want: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
		{
			name: "Test2",
			args: args{
				d: []string{"", "", ""},
			},
			want: "da39a3ee5e6b4b0d3255bfef95601890afd80709",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA1(tt.args.d...); got != tt.want {
				t.Errorf("SHA1() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestMD5(t *testing.T) {
	type args struct {
		d []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test1",
			args: args{
				d: []string{"aaa", "bbb"},
			},
			want: "6547436690a26a399603a7096e876a2d",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := MD5(tt.args.d...); got != tt.want {
				t.Errorf("MD5() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestSHA256(t *testing.T) {
	type args struct {
		d []string
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "Test1",
			args: args{
				d: []string{"aaa", "bbb"},
			},
			want: "2ce109e9d0faf820b2434e166297934e6177b65ab9951dbc3e204cad4689b39c",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := SHA256(tt.args.d...); got != tt.want {
				t.Errorf("SHA256() = %v, want %v", got, tt.want)
			}
		})
	}
}
