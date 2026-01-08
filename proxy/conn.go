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
	"sync/atomic"
	"time"

	_ "github.com/lbwise/learning/compressor/algos"
)

func NewConn(clientConn, destConn net.Conn, lg *log.Logger) *Conn {
	return &Conn{
		ID:         GenNextConnId(),
		ClientConn: clientConn,
		DestConn:   destConn,
		Logger:     lg,
		requests:   make([]*Request, 0),
		createdAt:  time.Now(),
	}
}

var (
	connIdCounter uint64
)

func GenNextConnId() uint64 {
	return atomic.AddUint64(&connIdCounter, 1)
}

type Conn struct {
	ID         uint64
	ClientConn net.Conn
	DestConn   net.Conn
	Logger     *log.Logger
	requests   []*Request
	createdAt  time.Time
	readSize   int
}

// We need to be able to accept multiple requests from the same connection
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

	req := NewRequest(raw)
	c.requests = append(c.requests, req)

	// Take incoming requests, mutate and forward
	httpReq, err := http.ReadRequest(bufio.NewReader(req.raw))
	if err != nil {
		c.Error(err)
		return
	}

	c.Log("[REQ-%d] %s %s host=%s agent=%s", req.ID, httpReq.Method, httpReq.URL, httpReq.Host, httpReq.UserAgent())

	n, err := io.Copy(c.DestConn, bytes.NewReader(raw))
	if err != nil {
		c.Error(err)
		return
	}

	c.Log("[REQ-%d] forwarded request of %d bytes from %s to %s ", req.ID, n, httpReq.Host, c.DestConn.RemoteAddr())

	// Take outgoing responses, mutate and forward
	c.DestConn.SetDeadline(time.Now().Add(5000 * time.Millisecond))

	body, err := ReadFromConn(c.DestConn)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = io.Copy(c.ClientConn, bytes.NewReader(body))
	if err != nil {
		c.Fatal(err)
		return
	}

	c.Log("[REQ-%d] connection to %s closed (%dms)", req.ID, c.ClientConn.RemoteAddr(), time.Now().Sub(start).Milliseconds())
}

func (c *Conn) Log(msg string, args ...interface{}) {
	c.Logger.Println(fmt.Sprintf("[CONN-%d]: %s", c.ID, fmt.Sprintf(msg, args...)))
}

func (c *Conn) Error(err error) {
	c.Logger.Println("[ERR]", err.Error())
}

func (c *Conn) Fatal(err error) {
	c.Logger.Println("[ERR]", err.Error())
}
