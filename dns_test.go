package gotik

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestClient_GetDNS(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if dns, err := c.GetDNS(); err != nil {
			t.Errorf("GetDNS() error = %v", err)
		} else {
			spew.Dump(dns)
		}
	})
}

func TestClient_SetDNS(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if dns, err := c.GetDNS(); err != nil {
			t.Errorf("GetDNS() error = %v", err)
		} else {
			spew.Dump(dns)
			dns.MaxUDPPacketSize = 1024
			dns.CacheMaxTTL = time.Hour * 24
			dns.Servers = []string{"8.8.8.8", "8.8.4.4"}
			dns.Servers = []string{}
			//dns.UseDOHServer = "https://dns.google/dns-query"
			dns.UseDOHServer = ""
			dns.VerifyDOHCert = false
			if err := c.SetDNS(dns); err != nil {
				t.Errorf("SetDNS() error = %v", err)
			} else {
				t.Logf("SetDNS() success")
			}
		}
	})
}
