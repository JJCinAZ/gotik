package main

import (
	tls2 "crypto/tls"
	"log"
	"os"
	"time"

	"github.com/jjcinaz/gotik"
)

func main() {
	var (
		routerConn *routeros.Client
		err        error
	)
	tls := tls2.Config{
		InsecureSkipVerify: true,
	}
	routerConn, err = routeros.DialTLSTimeout(os.Args[1], os.Args[2], os.Args[3], &tls, time.Second*10)
	if err != nil {
		log.Printf("unable to connect to router: %s", err)
		return
	}
	err = routerConn.RemovePPPSecretByName(os.Args[4])
	log.Println(err)

}
