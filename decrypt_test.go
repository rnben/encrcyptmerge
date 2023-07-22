package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestDecryptMap_MergeMap(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		dst string
		src []string
	}
	tests := []struct {
		name   string
		fields fields
		args   args
		want   map[string]interface{}
	}{
		{
			name:   "success",
			fields: fields{},
			args: args{
				dst: `{"a": "aa"}`,
				src: []string{},
			},
			want: map[string]interface{}{
				"a": "aa",
			},
		},
		{
			name:   "failed",
			fields: fields{},
			args: args{
				dst: `{"a" "aa"}`,
				src: []string{},
			},
			want: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DecryptMap{
				err: tt.fields.err,
			}
			if got := d.MergeMap(tt.args.dst, tt.args.src...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DecryptMap.MergeMap() = %v, want %v", got, tt.want)
			}
			_ = d.Err()
		})
	}
}

func TestDecryptMap_ProcessMap(t *testing.T) {
	type fields struct {
		fields string
		err    error
	}
	type args struct {
		data     map[string]interface{}
		filePath []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name:   "success",
			fields: fields{fields: "a"},
			args: args{
				data:     map[string]interface{}{"a": "aa"},
				filePath: []string{"a.log"},
			},
			wantErr: false,
		},
		{
			name:   "failed: merge",
			fields: fields{fields: "a", err: errors.New("merge err")},
			args: args{
				data:     map[string]interface{}{"a": 1},
				filePath: []string{"a.log"},
			},
			wantErr: true,
		},
		{
			name:   "failed: not string",
			fields: fields{fields: "a"},
			args: args{
				data:     map[string]interface{}{"a": 1},
				filePath: []string{"a.log"},
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &DecryptMap{
				err: tt.fields.err,
			}
			sensitiveFields = tt.fields.fields
			if err := d.ProcessMap(tt.args.data, tt.args.filePath...); (err != nil) != tt.wantErr {
				t.Errorf("DecryptMap.ProcessMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
