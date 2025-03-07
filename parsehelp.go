package gotik

import (
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"
)

var (
	reRosId         = regexp.MustCompile(`(?i)^\*[0-9A-F]{1,8}$`)
	reRosDuration   = regexp.MustCompile(`([+-]?\d+)([wdhmsun]{1,2})`)
	reRosTimeOffset = regexp.MustCompile(`([0-9.]+) ([wdhmsun]{1,2})`)
)

func isRosId(s string) bool {
	return reRosId.MatchString(s)
}

func parseTime(s string) time.Time {
	t, _ := time.ParseInLocation("Jan/02/2006 15:04:05", s, time.Local)
	return t
}

func parseBool(s string) bool {
	b := false
	if len(s) > 0 {
		switch s[0] {
		case 'y', 'Y', '1', 't', 'T':
			b = true
		}
	}
	return b
}

func splitFloat32(s string) [2]float32 {
	var x [2]float32
	if a := strings.SplitN(s, "/", 3); len(a) >= 2 {
		for i := 0; i < 2; i++ {
			if f, err := strconv.ParseFloat(a[i], 32); err == nil {
				x[i] = float32(f)
			}
		}
	}
	return x
}

func splitInt(s string) [2]int {
	var x [2]int
	if a := strings.SplitN(s, "/", 3); len(a) >= 2 {
		for i := 0; i < 2; i++ {
			if n, err := strconv.ParseInt(a[i], 10, 64); err == nil {
				x[i] = int(n)
			}
		}
	}
	return x
}

func splitString2(s string) [2]string {
	var x [2]string
	if a := strings.SplitN(s, "/", 3); len(a) >= 2 {
		for i := 0; i < 2; i++ {
			x[i] = a[i]
		}
	}
	return x
}

func parseInt(s string) int {
	i, _ := strconv.ParseInt(s, 10, 64)
	return int(i)
}

func parseHex(s string) int {
	var i int64
	if len(s) > 2 && s[0:2] == "0x" {
		i, _ = strconv.ParseInt(s[2:], 16, 64)
	} else {
		i, _ = strconv.ParseInt(s, 16, 64)
	}
	return int(i)
}

func parseFloat32(s string) float32 {
	f, _ := strconv.ParseFloat(s, 32)
	return float32(f)
}

func parseDuration(s string) time.Duration {
	var (
		nsecs  time.Duration
		negSet time.Duration
	)
	if len(s) == 0 {
		return nsecs
	}
	negSet = 1
	if s[0] == '+' {
		s = s[1:]
	} else if s[0] == '-' {
		s = s[1:]
		negSet = -1
	}
	if a := reRosTimeOffset.FindStringSubmatch(s); a != nil {
		n, err := strconv.ParseFloat(a[1], 64)
		if err != nil {
			return nsecs
		}
		switch a[2] {
		case "h":
			nsecs += 3600 * time.Second
		case "m":
			nsecs += 60 * time.Second
		case "s":
			nsecs += time.Second
		case "ms":
			nsecs += time.Millisecond
		case "us":
			nsecs += time.Microsecond
		case "ns":
			nsecs += time.Nanosecond
		}
		return time.Duration(n*float64(nsecs)) * negSet
	}
	a := reRosDuration.FindAllStringSubmatch(s, -1)
	if a != nil {
		for _, m := range a {
			// m[0] = "24w", m[1] = "24", m[2] = "w"
			n, err := strconv.ParseInt(m[1], 10, 64)
			if err != nil {
				break
			}
			switch m[2] {
			case "w":
				nsecs += time.Duration(n) * 86400 * 7 * time.Second
			case "d":
				nsecs += time.Duration(n) * 86400 * time.Second
			case "h":
				nsecs += time.Duration(n) * 3600 * time.Second
			case "m":
				nsecs += time.Duration(n) * 60 * time.Second
			case "s":
				nsecs += time.Duration(n) * time.Second
			case "ms":
				nsecs += time.Duration(n) * time.Millisecond
			case "us":
				nsecs += time.Duration(n) * time.Microsecond
			case "ns":
				nsecs += time.Duration(n) * time.Nanosecond
			}
		}
	}
	return nsecs * negSet
}

func parseTikObject(props map[string]string, i interface{}) {
	dstObj := reflect.ValueOf(i).Elem()
	dstObjType := dstObj.Type()
	for i := 0; i < dstObj.NumField(); i++ {
		field := dstObjType.Field(i)
		tag := field.Tag.Get("tik")
		if len(tag) > 0 {
			input, found := props[tag]
			if !found {
				continue
			}
			switch field.Type.Name() {
			case "bool":
				dstObj.Field(i).SetBool(parseBool(input))
			case "string":
				dstObj.Field(i).SetString(input)
			case "int":
				dstObj.Field(i).SetInt(int64(parseInt(input)))
			case "Duration":
				dstObj.Field(i).Set(reflect.ValueOf(parseDuration(input)))
			case "Time":
				if dstObj.Field(i).CanAddr() {
					dstObj.Field(i).Set(reflect.ValueOf(parseTime(input)))
				}
			}
		}
	}
}
