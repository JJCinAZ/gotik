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
	defer routerConn.Close()

	if err = routerConn.ExportConfig(); err != nil {
		log.Printf("unable to export config: %s", err)
		return
	}
}
