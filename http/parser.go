package http

import (
	_ "net/http"
)

type Request struct {
	Method string
}

func NewRequest(req []byte) (*Request, error) {
	httpReq := new(Request)
	httpReq.Method = string(req[:12])
	return httpReq, nil
}

//GET /some/path HTTP/1.1\r\n
//Host: example.com\r\n
//User-Agent: fake-client\r\n
//Accept: */*\r\n
//\r\n
