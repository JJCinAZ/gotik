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
		reply      []gotik.SNMPCommunity
		x          string
	)

	routerConn, err = gotik.DialTimeout(os.Args[1], os.Args[2], os.Args[3], time.Second*10)
	if err != nil {
		log.Printf("unable to connect to router: %s", err)
		return
	}
	reply, err = routerConn.GetSNMPCommunities()
	log.Println(reply, err)
	if err != nil {
		return
	}

	community := gotik.SNMPCommunity{
		Name:        "xxxx",
		WriteAccess: true,
		ReadAccess:  true,
	}
	x, err = routerConn.AddSNMPCommunity(community)
	log.Println(x, err)
	if err != nil {
		return
	}

	err = routerConn.RemoveSNMPCommunity("*0") // Should fail -- can't remove default community
	log.Println(err)

	reply, err = routerConn.GetSNMPCommunities()
	log.Println(reply, err)
	if err != nil {
		return
	}
	for i := range reply {
		if reply[i].Name == "xxxx" {
			err = routerConn.RemoveSNMPCommunity(reply[i].ID)
			log.Println(err)
			break
		}
	}

	reply, err = routerConn.GetSNMPCommunities()
	log.Println(reply, err)
}
