package common

import (
	"bufio"
	"encoding/gob"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"
)

type Address struct {
	Type    string
	City    string
	Country string
}

type VCard struct {
	FirstName string
	LastName  string
	Addresses []*Address
	Remark    string
}

func TestGob2File(t *testing.T) {
	pa := &Address{"private", "Aartselaar", "Belgium"}
	wa := &Address{"work", "Boom", "Belgium"}
	vc := VCard{"Jan", "Kersschot", []*Address{pa, wa}, "none"}
	// fmt.Printf("%v: \n", vc) // {Jan Kersschot [0x126d2b80 0x126d2be0] none}:
	// using an encoder:
	file, _ := os.OpenFile("../testdata/vcard.txt", os.O_CREATE|os.O_WRONLY, 0666)
	defer file.Close()
	enc := gob.NewEncoder(file)
	err := enc.Encode(vc)
	if err != nil {
		log.Println("Error in encoding gob")
	}
}

func TestFile2Gob(t *testing.T) {
	var vc VCard
	file, _ := os.Open("../testdata/vcard.txt")
	defer file.Close()
	inReader := bufio.NewReader(file)
	dec := gob.NewDecoder(inReader)
	err := dec.Decode(&vc)
	if err != nil {
		log.Printf("Error:%s in decoding gob", err)
	}

	fmt.Println(vc) //为什么Addresses 会是内存地址呢
	bs, err := json.Marshal(vc)
	fmt.Println(string(bs))
}
