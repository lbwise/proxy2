package proxy

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"time"

	_ "github.com/lbwise/learning/compressor/algos"
)

func NewConn(clientConn, destConn net.Conn, lg *log.Logger) *Conn {
	return &Conn{
		ID:         1,
		ClientConn: clientConn,
		DestConn:   destConn,
		Logger:     lg,
		requests:   make([]*Request, 0),
		createdAt:  time.Now(),
	}
}

type Conn struct {
	ID         int64
	ClientConn net.Conn
	DestConn   net.Conn
	Logger     *log.Logger
	requests   []*Request
	createdAt  time.Time
	readSize   int
}

func (c *Conn) Handle(ctx context.Context) {
	// handle error
	if c.DestConn == nil {
		return
	} else if c.ClientConn == nil {
		return
	}

	defer func() {
		c.DestConn.Close()
		c.ClientConn.Close()
	}()

	start := time.Now()

	raw, err := ReadFromConn(c.ClientConn)
	if err != nil {
		c.Error(err)
		return
	}

	fmt.Println("conn msg: ", string(raw))
	req := NewRequest(raw)
	c.requests = append(c.requests, req)

	// Take incoming requests, mutate and forward
	httpReq, err := http.ReadRequest(bufio.NewReader(req.raw))
	if err != nil {
		c.Error(err)
		return
	}

	c.Log("[REQ] %s %s host=%s agent=%s", httpReq.Method, httpReq.URL, httpReq.Host, httpReq.UserAgent())

	fmt.Println("THE INPUT")
	n, err := io.Copy(c.DestConn, bytes.NewReader(raw))
	if err != nil {
		c.Error(err)
		return
	}

	c.Log("forwarded request of %d bytes from %s to %s ", n, httpReq.Host, c.DestConn.RemoteAddr())

	// Take outgoing responses, mutate and forward
	c.DestConn.SetDeadline(time.Now().Add(5000 * time.Millisecond))

	body, err := ReadFromConn(c.DestConn)
	if err != nil {
		c.Error(err)
		return
	}
	c.Logger.Println(body)

	_, err = io.Copy(c.ClientConn, bytes.NewReader(body))
	if err != nil {
		c.Logger.Println("2 THIS ERROR HERE")
		c.Fatal(err)
		return
	}

	c.Log("Connection to %s closed (%dms)", c.ClientConn.RemoteAddr(), time.Now().Sub(start).Milliseconds())
}

func (c *Conn) Log(msg string, args ...interface{}) {
	c.Logger.Println(fmt.Sprintf("conn-%d: %s", c.ID, fmt.Sprintf(msg, args...)))
}

func (c *Conn) Error(err error) {
	c.Logger.Println("[ERR]", err.Error())
}

func (c *Conn) Fatal(err error) {
	c.Logger.Println("[ERR]", err.Error())
}
