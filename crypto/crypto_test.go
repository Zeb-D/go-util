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
	fmt.Println(t1, "-t1->", t1)
	t2 := AESDecryptByECB(t1, MyTestSecret)
	fmt.Println(t2, "-t2->", string(t2))
	data := "ebeced919187d1d0cbe1a5b1a19e54da950f31e3eac000e0a0e2ce14887efd6f393ca7311d039802d65625b70564c58b40e3f27f1149e66fbfca9d58140b1faa12ee0ad34b0e1b21c5e16274538ebbc679f52889ccf336df36532cc38d6fa3e997fde502c8691a669c1eb9244077bf9bf57ea9b0753c84e6cc2c06b7b39e88eda9be1a0744db82fe97cee42bc0ada4e5"
	secretKey := "b8bed00f89d562ae"
	t3 := AESDecryptByECB(data, secretKey)
	fmt.Println(t3, "-t3->", string(t3))
}

func TestDeCrypt22(t *testing.T) {
	data := "ebeced919187d1d0cbe1a5b1a19e54da950f31e3eac000e0a0e2ce14887efd6f393ca7311d039802d65625b70564c58b40e3f27f1149e66fbfca9d58140b1faa12ee0ad34b0e1b21c5e16274538ebbc679f52889ccf336df36532cc38d6fa3e997fde502c8691a669c1eb9244077bf9bf57ea9b0753c84e6cc2c06b7b39e88eda9be1a0744db82fe97cee42bc0ada4e5"
	secretKey := "b8bed00f89d562ae"
	d, e := AesECBDecrypt([]byte(data), []byte(secretKey))
	fmt.Println(e)
	fmt.Println(string(d))
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

func TestMD5(t *testing.T) {
	s := "1112221212"
	fmt.Println(MD5(s))
}

func TestAesECBEncrypt(t *testing.T) {
	orignStr := `{"biz_type":1,"data_type":1,"gw_id":"eb8530fb82fbb2b819jots","page_no":200,"page_size":200,"t":1639646802,"task_id":"","trans_id":""}`
	secretKey := "b8bed00f89d562ae"
	b, _ := AesECBEncrypt([]byte(orignStr), []byte(secretKey))
	fmt.Println(string(b))
	fmt.Println(base64.StdEncoding.EncodeToString(b))

	s := base64.StdEncoding.EncodeToString([]byte(orignStr))
	b, _ = AesECBEncrypt([]byte(s), []byte(secretKey))
	fmt.Println(string(b))
	fmt.Println(base64.StdEncoding.EncodeToString(b))

	orignStr = "{\"biz_type\":1,\"data_type\":1,\"gw_id\":\"eb8530fb82fbb2b819jots\",\"page_no\":200,\"page_size\":200,\"t\":1639646802,\"task_id\":\"\",\"trans_id\":\"\"}"
	b, _ = AesECBEncrypt([]byte(orignStr), []byte(secretKey))
	fmt.Println(string(b))
	fmt.Println(base64.StdEncoding.EncodeToString(b))

	s = base64.StdEncoding.EncodeToString([]byte(orignStr))
	b, _ = AesECBEncrypt([]byte(s), []byte(secretKey))
	fmt.Println(string(b))
	fmt.Println(base64.StdEncoding.EncodeToString(b))
}
