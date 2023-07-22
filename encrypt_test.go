package main

import (
	"errors"
	"reflect"
	"testing"
)

func TestEncryptMap_MergeMap(t *testing.T) {
	type fields struct {
		err error
	}
	type args struct {
		cur  string
		last []string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    map[string]interface{}
		wantErr bool
	}{
		{
			name:   "success0",
			fields: fields{},
			args: args{
				cur:  `{"a":"aa","b":"bb","c":"cc"}`,
				last: []string{`{"a":"a","b":"b","c":"c"}`},
			},
			want: map[string]interface{}{"a": "aa", "b": "bb", "c": "cc"},
		},
		{
			name:   "success1",
			fields: fields{},
			args: args{
				cur:  `{"a":"aa","b":"bb","c":"cc"}`,
				last: []string{`{"a":"a","b":"b"}`},
			},
			want: map[string]interface{}{"a": "aa", "b": "bb", "c": "cc"},
		},
		{
			name:   "success2",
			fields: fields{},
			args: args{
				cur:  `{"a":"aa","b":"bb","d":"dd"}`,
				last: []string{`{"a":"a","c":"c"}`},
			},
			want: map[string]interface{}{"a": "aa", "b": "bb", "d": "dd", "c": "c"},
		},
		{
			name:   "failed: src nil",
			fields: fields{},
			args: args{
				cur:  `{"a":"aa"}`,
				last: []string{},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "failed: dst invalid",
			fields: fields{},
			args: args{
				cur:  `{"a""aa"}`,
				last: []string{`{"a":"a"}`},
			},
			want:    nil,
			wantErr: true,
		},
		{
			name:   "failed: src invalid",
			fields: fields{},
			args: args{
				cur:  `{"a":"aa"}`,
				last: []string{`{"a""a"}`},
			},
			want:    nil,
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EncryptMap{
				err: tt.fields.err,
			}
			if got := e.MergeMap(tt.args.cur, tt.args.last...); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("EncryptMap.MergeMap() got = %v, want = %v", got, tt.want)
			}

			if err := e.Err(); err != nil && !tt.wantErr {
				t.Errorf("EncryptMap.MergeMap() err = %v, wantErr = %v", err, tt.wantErr)
			}
		})
	}
}

func TestEncryptMap_ProcessMap(t *testing.T) {
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
			name: "success1",
			args: args{
				data:     map[string]interface{}{"a": "aa"},
				filePath: []string{"s.log"},
			},
		},
		{
			name: "success2",
			args: args{
				data:     map[string]interface{}{"a": "aa"},
				filePath: []string{"s.log"},
			},
			fields: fields{
				fields: "a,z",
			},
		},
		{
			name: "failed: fielPath",
			args: args{
				data:     map[string]interface{}{"a": "aa"},
				filePath: []string{},
			},
			wantErr: true,
		},
		{
			name: "failed: merge err",
			args: args{
				data:     map[string]interface{}{"a": 1},
				filePath: []string{"a.log"},
			},
			fields: fields{
				err:    errors.New("merge error"),
				fields: "a",
			},
			wantErr: true,
		},
		{
			name: "failed: invalid str",
			args: args{
				data:     map[string]interface{}{"a": 1},
				filePath: []string{"a.log"},
			},
			fields:  fields{fields: "a"},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &EncryptMap{
				err: tt.fields.err,
			}
			sensitiveFields = tt.fields.fields
			if err := e.ProcessMap(tt.args.data, tt.args.filePath...); (err != nil) != tt.wantErr {
				t.Errorf("EncryptMap.ProcessMap() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
