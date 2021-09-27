package gotik

import (
	"testing"
	"time"
)

func Test_isRosId(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "empty test", args: args{s: ""}, want: false},
		{name: "missing prefix", args: args{s: "7"}, want: false},
		{name: "space", args: args{s: "* 34"}, want: false},
		{name: "non number", args: args{s: "*AG"}, want: false},
		{name: "too long", args: args{s: "*123456789"}, want: false},
		{name: "hex id", args: args{s: "*Ff783"}, want: true},
		{name: "decimal id", args: args{s: "*1"}, want: true},
		{name: "leading zeros", args: args{s: "*004d7"}, want: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isRosId(tt.args.s); got != tt.want {
				t.Errorf("isRosId() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseBool(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want bool
	}{
		{name: "true", args: args{s: "true"}, want: true},
		{name: "yes", args: args{s: "yes"}, want: true},
		{name: "false", args: args{s: "false"}, want: false},
		{name: "no", args: args{s: "no"}, want: false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseBool(tt.args.s); got != tt.want {
				t.Errorf("parseBool() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseFloat32(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want float32
	}{
		{name: "zero", args: args{s: "0"}, want: 0.0},
		{name: "negative", args: args{s: "-46.478"}, want: -46.478},
		{name: "positive", args: args{s: "983467.37674"}, want: 983467.37674},
		{name: "positive", args: args{s: "1"}, want: 1.0},
		{name: "positive", args: args{s: "+4773"}, want: 4773.0},
		{name: "empty", args: args{s: ""}, want: 0.0},
		{name: "garbage", args: args{s: "abc"}, want: 0.0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseFloat32(tt.args.s); got != tt.want {
				t.Errorf("parseFloat32() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_parseInt(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{name: "zero", args: args{s: "0"}, want: 0},
		{name: "negative", args: args{s: "-46"}, want: -46},
		{name: "positive", args: args{s: "983467"}, want: 983467},
		{name: "positive", args: args{s: "1"}, want: 1},
		{name: "positive", args: args{s: "+4773"}, want: 4773},
		{name: "empty", args: args{s: ""}, want: 0},
		{name: "garbage", args: args{s: "abc"}, want: 0},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseInt(tt.args.s); got != tt.want {
				t.Errorf("parseInt() = %v, want %v", got, tt.want)
			}
		})
	}
}

type testTikObject struct {
	ID           string        `tik:".id"`
	List         string        `tik:"list"`
	Disabled     bool          `tik:"disabled"`
	PktSize      int           `tik:"pkt-size"`
	CreationTime time.Time     `tik:"created"`
	Timeout      time.Duration `tik:"timeout"`
}

func makeTestTime(s string) time.Time {
	t, _ := time.ParseInLocation("Jan/02/2006 15:04:05", s, time.Local)
	return t
}

func Test_parseTikObject(t *testing.T) {
	type args struct {
		props map[string]string
	}
	tests := []struct {
		name string
		args args
		want testTikObject
	}{
		{
			name: "test1",
			args: args{
				props: map[string]string{
					".id":      "*1",
					"list":     "Management",
					"disabled": "false",
					"pkt-size": "1502",
					"created":  "feb/14/2044 14:30:30",
					"timeout":  "22d14h36m0s",
				},
			},
			want: testTikObject{
				ID:           "*1",
				List:         "Management",
				Disabled:     false,
				PktSize:      1502,
				CreationTime: makeTestTime("feb/14/2044 14:30:30"),
				Timeout:      1953360,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var got testTikObject
			parseTikObject(tt.args.props, &got)
			if got.Disabled != tt.want.Disabled ||
				got.Timeout != tt.want.Timeout ||
				got.List != tt.want.List ||
				got.ID != tt.want.ID ||
				got.PktSize != tt.want.PktSize ||
				got.CreationTime.Equal(tt.want.CreationTime) == false {
				t.Errorf("parseTikObject() = %v, want %v", got, tt.want)
			}
		})
	}
}

type parsedTime struct {
	Year                 int
	Month                time.Month
	Day                  int
	Hour, Minute, Second int // 15:04:05 is 15, 4, 5.
}

func sameTime(t time.Time, u parsedTime) bool {
	// Check aggregates.
	year, month, day := t.Date()
	hour, min, sec := t.Clock()
	if year != u.Year || month != u.Month || day != u.Day ||
		hour != u.Hour || min != u.Minute || sec != u.Second {
		return false
	}
	// Check individual entries.
	return t.Year() == u.Year &&
		t.Month() == u.Month &&
		t.Day() == u.Day &&
		t.Hour() == u.Hour &&
		t.Minute() == u.Minute &&
		t.Second() == u.Second
}

func Test_parseTime(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want parsedTime
	}{
		{
			name: "good 1", args: args{s: "oct/18/2006 16:24:41"},
			want: parsedTime{2006, time.October, 18, 16, 24, 41},
		},
		{
			name: "good 2", args: args{s: "dec/25/2029 09:00:01"},
			want: parsedTime{2029, time.December, 25, 9, 0, 1},
		},
		{
			name: "beyond 2038", args: args{s: "jan/10/2042 00:00:00"},
			want: parsedTime{2042, time.January, 10, 0, 0, 0},
		},
		{
			name: "missing leading zeros", args: args{s: "feb/14/2019 8:0:0"},
			want: parsedTime{1, time.January, 1, 0, 0, 0},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseTime(tt.args.s); !sameTime(got, tt.want) {
				t.Errorf("parseTime() = %v, want %v", got, tt.want)
			}
		})
	}
}
