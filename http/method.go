package http

import "fmt"

type Method int

func (m Method) String() string {
	switch m {
	case GET:
		return "GET"
	case POST:
		return "POST"
	case PUT:
		return "POST"
	case PATCH:
		return "POST"
	default:
		return "POST"
	}
}

const (
	GET Method = iota
	POST
	PUT
	PATCH
	DELETE
)

type Parser struct {
	msg []byte
}

func (p *Parser) Method(req *Request, method string) error {
	switch method {
	case "GET":
		req.Method = GET
		return nil
	case "POST":
		req.Method = POST
		return nil
	case "PATCH":
		req.Method = PATCH
		return nil
	case "PUT":
		req.Method = PUT
		return nil
	case "DELETE":
		req.Method = DELETE
		return nil
	default:
		return fmt.Errorf("invalid http method: %s", method)
	}
}
