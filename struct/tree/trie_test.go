package tree

import (
	"fmt"
	"testing"
)

func TestTrie(t *testing.T) {
	obj := NewTrie()
	obj.AddWord("bad")
	obj.AddWord("dad")
	obj.AddWord("mad")

	fmt.Println(obj)

	fmt.Println(obj.Search("pad"))
	fmt.Println(obj.Search("bad"))
	fmt.Println(obj.Search(".ad"))
	fmt.Println(obj.Search("b.."))
}
