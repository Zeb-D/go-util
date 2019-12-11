package set

import (
	"bytes"
	"fmt"
	"path/filepath"
	"testing"
	"time"

	"github.com/Zeb-D/go-util/common"
)

func testSet(set ISet, filename string) time.Duration {
	startTime := time.Now()

	words := common.ReadFile(filename)
	fmt.Println("Total words:", len(words), words)
	for _, word := range words {
		set.Add(word)
	}
	fmt.Println("Total different words:", set.Size())

	return time.Now().Sub(startTime)
}

func TestSet(t *testing.T) {
	filename, _ := filepath.Abs("../../testdata/consumer_config.yml")

	bstSet := NewBSTSet()
	time1 := testSet(bstSet, filename)
	fmt.Println("BST set :", time1)

	linkedListSet := NewLinkedListSet()
	time2 := testSet(linkedListSet, filename)
	fmt.Println("linkedList set:", time2)
}

func TestMorseRepresentations(t *testing.T) {
	words := []string{"gin", "zen", "gig", "msg"}

	fmt.Println(uniqueMorseRepresentations(words))
}

// Go 中没有 set 类型，这里使用 map 实现
func uniqueMorseRepresentations(words []string) int {
	morseCodes := []string{".-", "-...", "-.-.", "-..", ".", "..-.", "--.", "....", "..", ".---", "-.-", ".-..", "--", "-.", "---", ".--.", "--.-", ".-.", "...", "-", "..-", "...-", ".--", "-..-", "-.--", "--.."}

	buffer := bytes.Buffer{}
	uniqueWord := make(map[string]bool)
	for _, word := range words {
		buffer.Reset()
		for _, letter := range word {
			//fmt.Println(letter, "->", letter-'a', "->", morseCodes[letter-'a'])
			buffer.WriteString(morseCodes[letter-'a'])
			//fmt.Println("buffer1->", buffer)
		}
		//fmt.Println("buffer2->", buffer.String())
		uniqueWord[buffer.String()] = true
	}

	return len(uniqueWord)
}
