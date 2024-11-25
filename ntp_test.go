package gotik

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestClient_GetNTPClient(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if ntp, err := c.GetNTPClient(); err != nil {
			t.Errorf("GetNTPClient() error = %v", err)
		} else {
			spew.Dump(ntp)
		}
	})
}

func TestClient_SetNTPClient(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if ntp, err := c.GetNTPClient(); err != nil {
			t.Errorf("GetDNS() error = %v", err)
		} else {
			switch n := ntp.(type) {
			case NTPClient6:
				n.Servers = []string{"64.119.32.60", "64.119.32.6"}
				n.Enabled = true
				ntp = n
			case NTPClient7:
				n.Servers = []string{"time.google.com", "time.cloudflare.com"}
				n.Enabled = true
				ntp = n
			}
			if err := c.SetNTPClient(ntp); err != nil {
				t.Errorf("SetNTPClient() error = %v", err)
			} else {
				t.Logf("SetNTPClient() success")
			}
		}
	})
}
