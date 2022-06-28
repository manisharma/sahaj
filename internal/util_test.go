package internal

import (
	"reflect"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestNow(t *testing.T) {
	now := time.Now()
	tests := []struct {
		name string
		want time.Time
	}{
		{
			name: "should set to today's start time",
			want: time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, now.Location()),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := Now(); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("Now() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestToStringPtr(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want *string
	}{
		{
			name: "should return pointer",
			args: args{
				s: "",
			},
			want: ToStringPtr(""),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := ToStringPtr(tt.args.s)
			assert.NotNil(t, got, "result must not be nil")
		})
	}
}
