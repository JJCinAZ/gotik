package gotik

import (
	"os"
	"testing"
	"time"
)

func TestClient_AddFile(t *testing.T) {
	c, err := DialTimeout(os.Getenv("RTR"), os.Getenv("RTRUSER"), os.Getenv("RTRPASS"),
		time.Second*5)
	if err != nil {
		t.Fatalf("connection failure %v", err)
	}
	t.Run("test1", func(t *testing.T) {
		if err := c.AddFile("test.txt", "test contents here"); err != nil {
			t.Errorf("AddFile() error = %v", err)
		}
	})
}
