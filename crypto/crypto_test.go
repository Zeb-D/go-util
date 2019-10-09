package crypto

import (
	"crypto/aes"
	"encoding/base64"
	"fmt"
	"strings"
	"testing"
)

const MyTestSecret = "12C21o31m34mo2na"

func TestSign(t *testing.T) {
	var body = "my name is go-util"
	var time = "1568982719748"
	var secret = MyTestSecret
	var s = fmt.Sprintf("t=%s&data=%s&s=%s", time, body, secret)
	fmt.Printf("before:%s, after:%s,length:%d \n", s,
		Sha256Sign(s),
		len(Sha256Sign(s)))
}

func TestAESByECB(t *testing.T) {
	var text = "aa1ba24"
	t1 := AESEncryptByECB(text, MyTestSecret)
	fmt.Println(t1, "-->", t1)
	t2 := AESDecryptByECB(t1, MyTestSecret)
	fmt.Println(t2, "-->", string(t2))
}

func TestDecryptAES(t *testing.T) {
	rb := []byte(MyTestSecret)
	b := make([]byte, 16)
	strings.NewReader("aaa").Read(b)
	// b=b[0:16];
	fmt.Print("b:", b)
	cip, _ := aes.NewCipher(b)
	fmt.Print("cip:", cip, "err:")
	out := make([]byte, len(rb))
	cip.Encrypt(rb, out)
	fmt.Println()
	fmt.Println(base64.StdEncoding.EncodeToString(rb))
}
