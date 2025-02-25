package gotik

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestClient_GetSNMP(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if s, err := c.GetSNMP(); err != nil {
			t.Errorf("GetSNMP() error = %v", err)
		} else {
			spew.Dump(s)
		}
	})
}

func TestClient_SetSNMP(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if s, err := c.GetSNMP(); err != nil {
			t.Errorf("GetSNMP() error = %v", err)
		} else {
			spew.Dump(s)
			s.Enabled = false
			s.Contact = ""
			s.TrapTarget = ""
			s.TrapGenerators = []string{"start-trap"}
			if err := c.SetSNMP(s); err != nil {
				t.Errorf("SetSNMP() error = %v", err)
			} else {
				t.Logf("SetSNMP() success")
			}
		}
	})
}
