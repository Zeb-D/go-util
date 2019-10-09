package interact

import (
	"fmt"
	"testing"
	"time"
)

func TestTicker(t *testing.T) {
	ticker := time.NewTicker(100 * time.Millisecond)
	go func() {
		time.Sleep(2e9)
		ticker.Stop()
	}()
	for {
		select {
		case t := <-ticker.C:

			fmt.Println(t)
		case <-time.After(3e9):
			fmt.Println("break")
			return
		}
	}
}
