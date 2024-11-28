package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jjcinaz/gotik"
)

func main() {
	var (
		routerConn *gotik.Client
		err        error
		list       []gotik.IPv4FilterRule
	)
	routerConn, err = gotik.DialTimeout(os.Args[1], os.Args[2], os.Args[3], time.Second*10)
	if err != nil {
		log.Printf("unable to connect to router: %s", err)
		return
	}
	list, err = routerConn.GetIPv4Filters(os.Args[4])
	if err != nil {
		log.Println(err)
	} else {
		for _, item := range list {
			fmt.Printf("%#v\n", item)
		}
	}
}
