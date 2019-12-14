package Map

import (
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Zeb-D/go-util/common"
)

func TestMap(t *testing.T) {
	filename, _ := filepath.Abs("../../testdata/consumer_config.yml")

	bstMap := NewBSTMap()
	time1 := testMap(bstMap, filename)
	fmt.Println("BST map :", time1)

	linkedListMap := NewLinkedListMap()
	time2 := testMap(linkedListMap, filename)
	fmt.Println("linkedListMap set:", time2)
}

func testMap(p IMap, filename string) time.Duration {
	startTime := time.Now()

	words := common.ReadFile(filename)
	fmt.Println("Total words:", len(words))

	for _, word := range words {
		if p.Contains(word) {
			p.Set(word, p.Get(word).(int)+1)
		} else {
			p.Add(word, 1)
		}
	}
	fmt.Println("Total different words: ", p.Size())
	fmt.Println("Frequency of PRIDE:", p.Get("pride"))
	fmt.Println("Frequency of PREJUDICE: ", p.Get("prejudice"))

	return time.Now().Sub(startTime)
}
