package gotik

import (
	"fmt"
	"os"
	"testing"
)

func getRouterConn(t *testing.T) *Client {
	routerConn, err := Dial(os.Getenv("RTRIP"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"))
	if err != nil {
		t.Errorf("Dial error %s", err)
		return nil
	}
	return routerConn
}

func TestClient_GetUpdateInfo(t *testing.T) {
	routerConn := getRouterConn(t)
	defer routerConn.Close()
	info, _ := routerConn.GetUpdateInfo()
	fmt.Printf("%#v\n", info)
}

func TestClient_SetUpdateChannel(t *testing.T) {
	routerConn := getRouterConn(t)
	defer routerConn.Close()
	err := routerConn.SetUpdateChannel("long-term")
	fmt.Println(err)
}

func TestClient_DownloadUpdate(t *testing.T) {
	routerConn := getRouterConn(t)
	defer routerConn.Close()
	info, err := routerConn.DownloadUpdates()
	fmt.Println(err)
	fmt.Printf("%#v\n", info)
}

func TestClient_InstallUpdate(t *testing.T) {
	routerConn := getRouterConn(t)
	defer routerConn.Close()
	info, err := routerConn.InstallUpdates()
	fmt.Println(err)
	fmt.Printf("%#v\n", info)
}
