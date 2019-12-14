package union_find

import (
	"fmt"
	"math/rand"
	"testing"
	"time"
)

func TestUnionFind(t *testing.T) {
	size := 10000000
	m := 10000000
	uf6 := NewUnionFind(size)
	fmt.Println(uf6)
	fmt.Println(testUF(uf6, m))
}

func testUF(uf UF, m int) time.Duration {
	size := uf.Size()
	rand.Seed(time.Now().Unix())

	startTime := time.Now()

	for i := 0; i < m; i++ {
		a := rand.Intn(size)
		b := rand.Intn(size)
		uf.UnionElements(a, b)
	}

	for i := 0; i < m; i++ {
		a := rand.Intn(size)
		b := rand.Intn(size)
		uf.IsConnected(a, b)
	}

	return time.Now().Sub(startTime)
}
