package gotik

import (
	"fmt"
	"testing"
)

func TestClient_GetOspf2LsaTable(t *testing.T) {
	routerConn := getRouterConn(t)
	defer routerConn.Close()
	table, _ := routerConn.GetOspf2LsaTable()
	for _, e := range table {
		if e.LSAID == "64.119.44.124" {
			fmt.Printf("%#v\n", e)
		}
	}
}

func Test_parseNetworkMaskToBits(t *testing.T) {
	type args struct {
		s string
	}
	tests := []struct {
		name string
		args args
		want int
	}{
		{"test1", args{"255.255.0.0"}, 16},
		{"test2", args{"0.0.0.0"}, 0},
		{"test3", args{"255.255.224.0"}, 19},
		{"test4", args{"255.255.255.0"}, 24},
		{"test5", args{"255.255.255.240"}, 28},
		{"test6", args{"255.255.255.255"}, 32},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := parseNetworkMaskToBits(tt.args.s); got != tt.want {
				t.Errorf("parseNetworkMaskToBits() = %v, want %v", got, tt.want)
			}
		})
	}
}
