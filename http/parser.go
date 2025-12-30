package http

import (
	"fmt"
	_ "net/http"
	"strings"
)

//GET /some/path HTTP/1.1\r\n
//Host: example.com\r\n
//User-Agent: fake-client\r\n
//Accept: */*\r\n
//\r\n

func NewParser(msg []byte) *Parser {
	return &Parser{msg}
}

type Request struct {
	Method
	URI       string
	UserAgent string
	Host      string
}

func (req *Request) String() string {
	return fmt.Sprintf("Method: %s\nURI: %s\nHOST: %s\nUA: %s", req.Method, req.URI, req.UserAgent, req.Host)
}

// ONLY TAKES VALID REQUESTS RN LOL
func (p *Parser) Parse() (*Request, error) {
	httpReq := new(Request)

	lines := strings.Split(string(p.msg), "\n")

	line1 := strings.Split(lines[0], " ")
	err := p.Method(httpReq, line1[0])
	if err != nil {
		return nil, err
	}

	httpReq.URI = line1[1]
	httpReq.Host = strings.Split(lines[1], " ")[1]
	httpReq.UserAgent = strings.Split(lines[2], " ")[1]
	return httpReq, nil
}
