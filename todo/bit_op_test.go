package todo

import (
	"encoding/base32"
	"fmt"
	"github.com/Zeb-D/go-util/common"
	"github.com/Zeb-D/go-util/crypto/resolver"
	"testing"
)

//	位运算
func TestBitOp(t *testing.T) {
	fmt.Println(8 ^ 255)
	fmt.Println(1>>8 ^ 255)
}

func TestCheck(t *testing.T) {
	var hex string = "00000018053111EEF843174B9E8EB7041BA47A9338B5441C44E0BBBDD13D2392048643426C658091F0F5D32009E08982D215C3817AC1EDC478452D2A995F84A0CDADB799CF2C077D732773A372FDCC6E6089FA018BE012F33B2E9D884791D5DB92C9CCE0027041D0"
	fmt.Println(hex)
	bytes, err := base32.HexEncoding.DecodeString(hex)
	if err != nil {
		fmt.Println("err ->", err)
	}
	fmt.Println(bytes)
	b, _ := common.ToInt8s(hex)
	bs := resolver.Checksum(b[4:])
	fmt.Println("Checksum ->", bs)
	i := resolver.FromBytes(bs[0], bs[1], bs[2], bs[3])
	fmt.Println(i)
	//bs, _ := common.ToBytes(hex)
	//fmt.Println(bs[4:])
}
