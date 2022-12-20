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
		list       []gotik.Interface
	)
	routerConn, err = gotik.DialTimeout(os.Args[1], os.Args[2], os.Args[3], time.Second*10)
	if err != nil {
		log.Printf("unable to connect to router: %s", err)
		return
	}
	list, err = routerConn.GetInterfacesOfTypes(os.Args[4])
	if err != nil {
		log.Println(err)
	} else {
		for _, iface := range list {
			fmt.Println(iface)
		}
	}
}
