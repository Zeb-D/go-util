package http

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/magiconair/properties/assert"
	"io/ioutil"
	"net/http"
	"reflect"
	"testing"
)

type user struct {
	Name    string `json:"name"`
	Age     uint   `json:"age"`
	Sex     byte   `json:"sex"`
	Address string `json:"addr"`
}

func TestReadRequestBody(t *testing.T) {
	want := []byte("hello")
	args := struct {
		req  *http.Request
		copy bool
	}{
		req: &http.Request{
			Body: ioutil.NopCloser(bytes.NewReader(want)),
		},
		copy: true,
	}
	//test for request
	bs := ExtractRequestBody(args.req, args.copy)
	fmt.Println("request body1:", string(bs))

	bs1 := ExtractRequestBody(args.req, args.copy)
	fmt.Println("request body2:", string(bs1))
	assert.Equal(t, bs, want)

	//test for request by user
	user1 := user{
		Name:    "yda",
		Address: "SZ nanshan",
		Sex:     byte(1),
		Age:     16,
	}
	requestBody, _ := json.Marshal(user1)
	args.req.Body = ioutil.NopCloser(bytes.NewReader(requestBody))

	t.Run(user1.Name, func(t *testing.T) {
		bs := ExtractRequestBody(args.req, args.copy)
		println(string(bs))
		copyUser := &user{}
		json.Unmarshal(bs, copyUser)
		fmt.Println("copy copyUser:", copyUser)
		reflect.DeepEqual(copyUser, user1)
	})
}
