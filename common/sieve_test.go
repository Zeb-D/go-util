package common

import (
	"fmt"
	"testing"
)

func TestSieve(t *testing.T) {
	primes := Sieve()
	fmt.Println(<-primes)
	for prime := <-primes; prime < 100; prime = <-primes {
		fmt.Println(prime)
	}
}
