package http

import (
	"bytes"
	"io/ioutil"
	"net/http"
)

//关于http Request还是Response的body 流只能读取一次，那么如何做到拦截式输出到日志
func ExtractResponseBody(resp *http.Response, revert bool) []byte {
	if resp == nil || resp.Body == nil {
		return nil
	}
	bs, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil
	}

	if revert {
		newBody := make([]byte, len(bs))
		copy(newBody, bs)
		resp.Body = ioutil.NopCloser(bytes.NewReader(newBody))
	}
	return bs
}

func ExtractRequestBody(req *http.Request, revert bool) []byte {
	if req == nil || req.Body == nil {
		return nil
	}
	bs, err := ioutil.ReadAll(req.Body)
	defer req.Body.Close()
	if err != nil {
		return nil
	}

	if revert {
		newBody := make([]byte, len(bs))
		copy(newBody, bs)
		req.Body = ioutil.NopCloser(bytes.NewReader(newBody))
	}
	return bs
}
