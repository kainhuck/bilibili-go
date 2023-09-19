package utils

import (
	"testing"
	"time"
)

func TestGetCorrespondPath(t *testing.T) {
	type args struct {
		ts int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			"t1",
			args{ts: time.Now().UnixMilli()},
			false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetCorrespondPath(tt.args.ts)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetCorrespondPath() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			t.Logf("GetCorrespondPath() got = %v", got)
		})
	}
}
