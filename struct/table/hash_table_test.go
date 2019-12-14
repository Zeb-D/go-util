package table

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Zeb-D/go-util/common"
)

func TestHashTable(t *testing.T) {
	fmt.Println("consumer_config")

	filename, _ := filepath.Abs("../../testdata/consumer_config.yml")
	words := common.ReadFile(filename)
	fmt.Println("Total words: ", len(words))

	// Test RBTree
	startTime := time.Now()
	ht := NewHashTable()
	for _, word := range words {
		if ht.Contains(word) {
			ht.Set(word, ht.Get(word).(int)+1)
		} else {
			ht.Add(word, 1)
		}
	}
	for _, word := range words {
		ht.Contains(word)
	}
	fmt.Println("HashTable: ", time.Now().Sub(startTime))
}
