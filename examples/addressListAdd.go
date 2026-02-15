package main

import (
	"log"
	"os"
	"time"

	"github.com/jjcinaz/gotik"
)

func main() {
	var (
		routerConn *gotik.Client
		err        error
	)

	routerConn, err = gotik.DialTimeout(os.Args[1], os.Args[2], os.Args[3], time.Second*10)
	if err != nil {
		log.Printf("unable to connect to router: %s", err)
		return
	}
	reply, err := routerConn.AddIPv4AddressListEntry("badguys", "192.0.2.4", 60*time.Second,
		"xxxx", false)
	if err != nil {
		log.Println(err)
	} else {
		log.Println(reply)
	}
}
