package main

import (
	"reflect"
	"testing"
)

func TestNewJson(t *testing.T) {
	type args struct {
		action Action
	}
	tests := []struct {
		name string
		args args
		want string
	}{
		{
			name: "decrypt",
			args: args{
				action: ActionDecrypt,
			},
			want: "main.DecryptMap",
		},
		{
			name: "encrypt",
			args: args{
				action: ActionEncrypt,
			},
			want: "main.EncryptMap",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewJson(tt.args.action)
			if rt := reflect.TypeOf(got).Elem().String(); rt != tt.want {
				t.Errorf("NewJson() = %v, want %v", rt, tt.want)
			}

		})
	}
}
