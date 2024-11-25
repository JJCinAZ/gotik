package gotik

import (
	"os"
	"testing"
	"time"

	"github.com/davecgh/go-spew/spew"
)

func TestClient_GetQueueTree(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if list, err := c.GetQueueTree("non-exist"); err != nil {
			t.Errorf("GetQueueTree() error = %v", err)
		} else {
			spew.Dump(list)
		}
	})
}
