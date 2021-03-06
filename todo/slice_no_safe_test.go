package todo

import (
	"fmt"
	"sync"
	"testing"
)

func TestSlice(t *testing.T) {
	var (
		slc = []int{}
		n   = 10000
		wg  sync.WaitGroup
	)

	wg.Add(n)

	for i := 0; i < n; i++ {
		go func() {
			slc = append(slc, i)
			wg.Done()
		}()
	}

	wg.Wait()

	fmt.Println("len:", len(slc))
	fmt.Println("done")
}

func TestSliceLock(t *testing.T) {
	var (
		n   = 10000
		slc = make([]int, 0, n)

		wg   sync.WaitGroup
		lock sync.Mutex
	)

	for i := 0; i < n; i++ {
		wg.Add(1)
		go func(a int) {
			defer wg.Done()
			lock.Lock()
			defer lock.Unlock()
			slc = append(slc, a)
		}(i)
		wg.Wait()
	}

	fmt.Println("len:", len(slc))
}
