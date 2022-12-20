package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	scp "github.com/bramvdbogaerde/go-scp"
	"github.com/bramvdbogaerde/go-scp/auth"
	"github.com/jjcinaz/gotik"
	"golang.org/x/crypto/ssh"
)

func main() {
	type result struct {
		worked bool
		ip     string
		id     string
	}
	var (
		wg         sync.WaitGroup
		mtxResults sync.Mutex
	)
	user := flag.String("u", os.Getenv("RTRUSER"), "router user (will default to value of RTRUSER environment variable)")
	pass := flag.String("p", os.Getenv("RTRPASSWORD"), "router password (will default to value of RTRPASSWORD environment variable)")
	nobackup := flag.Bool("nobackup", false, "skips backup of router configuration before upgrade")
	configpath := flag.String("path", ".", "path into which configuration backups should be stored")
	mode := flag.String("m", "download", "mode: download = download new version only; install = download and reboot to install")
	channel := flag.String("c", "long-term", "upgrade channel to use (long-term, stable, testing)")
	flag.Parse()
	if flag.NArg() < 1 {
		fmt.Println("must supply IP address or DNS name of router(s) to upgrade")
		os.Exit(1)
	}
	routerList := flag.Args()
	wg.Add(len(routerList))
	results := make([]result, 0, len(routerList))
	for _, rtr := range routerList {
		thisRtr := rtr
		go func() {
			r := result{ip: thisRtr}
			r.id, r.worked = doUpgrade(thisRtr, *user, *pass, *configpath, *nobackup, *mode, *channel)
			mtxResults.Lock()
			results = append(results, r)
			mtxResults.Unlock()
			wg.Done()
		}()
	}
	wg.Wait()
	exitVal := 0
	for _, r := range results {
		if !r.worked {
			exitVal = 2
		}
	}
	os.Exit(exitVal)
}

func doUpgrade(rtr, user, pass string, configpath string, nobackup bool, mode string, channel string) (string, bool) {
	var (
		routerConn *gotik.Client
		err        error
		info       gotik.PackageUpdate
		id         string
	)
	routerConn, err = gotik.DialTimeout(rtr, user, pass, time.Second*5)
	if err != nil {
		log.Printf("%s: unable to connect to router: %s", rtr, err)
		return id, false
	}
	defer routerConn.Close()
	if id, err = routerConn.GetSystemId(); err != nil {
		log.Printf("%s: Unable to get system ID err: %s", rtr, err)
		return id, false
	}
	log.Printf("%s: connected to router %s", rtr, id)
	if !nobackup {
		if err = enableSSH(routerConn); err != nil {
			log.Printf("%s: error enabling SSH service: %s", rtr, err)
			return id, false
		}
		log.Printf("%s: creating configuration backup on router", rtr)
		configname := exportName(id)
		if err = routerConn.ExportConfig("/", configname, false); err != nil {
			log.Printf("%s: unable to export config: %s", rtr, err)
			return id, false
		}
		log.Printf("%s: downloading configuration to local path %s", rtr, configpath)
		if err = getConfig(rtr, user, pass, configname, filepath.Join(configpath, configname)); err != nil {
			log.Printf("%s: SCP config: %s", rtr, err)
			return id, false
		}
	}
	if err = routerConn.SetUpdateChannel(channel); err != nil {
		log.Printf("%s: unable to set channel: %s", rtr, err)
		return id, false
	}
	switch mode {
	case "download":
		if info, err = routerConn.DownloadUpdates(); err != nil {
			log.Printf("%s: unable to install updates: %s", rtr, err)
		} else {
			log.Printf("%s: %s", rtr, info.Status)
		}
	case "install":
		if info, err = routerConn.InstallUpdates(); err != nil {
			log.Printf("%s: unable to install updates: %s", rtr, err)
		} else {
			log.Printf("%s: %s", rtr, info.Status)
		}
	default:
		log.Printf("invalid mode, must be install or download")
		return id, false
	}
	return id, true
}

func enableSSH(routerConn *gotik.Client) error {
	list, err := routerConn.GetIPServices()
	if err != nil {
		return err
	}
	for i := range list {
		if list[i].Name == "ssh" && list[i].Disabled {
			err = routerConn.SetIPServiceDisable(list[i].ID, false)
			break
		}
	}
	return err
}

func exportName(routerid string) string {
	// Remove chars which might be problematic for file system
	routerid = strings.Map(func(r rune) rune {
		switch r {
		case '\\', '/', ':', '&', ' ':
			return '_'
		case '`', '"', '\'':
			return -1
		}
		return r
	}, routerid)
	return routerid + time.Now().Format("_20060102.rsc")
}

func getConfig(rtr, user, pass string, remotefile, localfile string) error {
	clientConfig, _ := auth.PasswordKey(user, pass, ssh.InsecureIgnoreHostKey())
	client := scp.NewClient(rtr+":22", &clientConfig)
	if err := client.Connect(); err != nil {
		return fmt.Errorf("failed to establish SSH connection: %s", err)
	}
	defer client.Close()

	f, err := os.Create(localfile)
	if err != nil {
		return fmt.Errorf("unable to create local file %s: %s", localfile, err)
	}
	defer f.Close()
	return client.CopyFromRemote(f, remotefile)
}
