package gotik

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestClient_GetPPPoEServers(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Errorf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if list, err := c.GetPPPoEServers("e2-v310-RogerTree"); err != nil {
			t.Errorf("GetPPPoEServers() error = %v", err)
		} else {
			spew.Dump(list)
		}
	})
	t.Run("test2", func(t *testing.T) {
		if list, err := c.GetPPPoEServers(""); err != nil {
			t.Errorf("GetPPPoEServers() error = %v", err)
		} else {
			spew.Dump(list)
		}
	})
}
