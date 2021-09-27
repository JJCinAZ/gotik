package main

import (
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

	routerConn, err = routeros.DialTimeout(os.Args[1], os.Args[2], os.Args[3], time.Second*10)
	if err != nil {
		log.Printf("unable to connect to router: %s", err)
		return
	}
	secret := routeros.PPPSecret{
		Name:     "1234",
		Password: "passwordhere",
		Service:  "sstp",
	}
	reply, err := routerConn.AddPPPSecret(secret)
	log.Println(reply, err)
}
