package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net"
	"net/http"
)

var (
	certFile = flag.String("certFile", "./todo/main/http/main/cert.pem", "certFile path")
	keyFile  = flag.String("keyFile", "./todo/main/http/main/key.pem", "keyFile path")
)

//1. 生成私钥
//openssl genrsa -out key.pem 2048
//2. 生成证书
//openssl req -new -x509 -key key.pem -out cert.pem -days 1095

type TestHello struct {
	A  string `json:"a"`
	B  string `json:"b"`
	ZZ string `json:"zz"`
}

func main() {
	log.Println("start: 443")
	l, err := net.Listen("tcp", ":443")
	if err != nil {
		log.Fatal(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/info", info)

	err = http.ServeTLS(l, mux, *certFile, *keyFile)
	if err != nil {
		log.Fatal(err)
	}
}

func info(w http.ResponseWriter, r *http.Request) {
	log.Println(r.Method, r.UserAgent())

	buf, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("recv:", string(buf))

	t := &TestHello{
		A:  "第一个字段",
		B:  "第二个字段",
		ZZ: "一个字段",
	}
	resp, _ := json.Marshal(t)
	_, err = w.Write(resp)
	if err != nil {
		log.Println(err)
	}
}
