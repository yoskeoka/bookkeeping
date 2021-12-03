package main

import (
	"flag"
	"testing"
	"time"
)

func Test_dateFlag(t *testing.T) {

	tests := []struct {
		name    string
		args    []string
		want    time.Time
		wantErr bool
	}{
		{
			name:    "no parameter",
			args:    []string{},
			want:    time.Time{},
			wantErr: false,
		},
		{
			name:    "valid date",
			args:    []string{"-d", "20060102"},
			want:    time.Date(2006, 1, 2, 0, 0, 0, 0, time.UTC),
			wantErr: false,
		},
		{
			name:    "wrong format date",
			args:    []string{"-d", "20060199"},
			want:    time.Time{},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			fset := flag.NewFlagSet("bk", flag.ContinueOnError)
			var d time.Time
			fset.Var(&dateFlag{&d}, "d", "date (format: yyyymmdd)")

			err := fset.Parse(tt.args)
			if (err != nil) != tt.wantErr {
				t.Errorf("flag.Parse() returns %v, but wantErr = %v", err, tt.wantErr)
			}
			if err != nil {
				return
			}

			if !d.Equal(tt.want) {
				t.Errorf("want %v, but got %v", tt.want, d)
			}
		})
	}

}
