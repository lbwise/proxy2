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

func NewConnHandler(clientConn, destConn net.Conn, lg *log.Logger) *ConnHandler {
	return &ConnHandler{
		ID:         GenNextConnId(),
		clientConn: clientConn,
		destConn:   destConn,
		logger:     lg,
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

type ConnHandler struct {
	ID         uint64
	clientConn net.Conn
	destConn   net.Conn
	logger     *log.Logger
	requests   []*Request
	createdAt  time.Time
	readSize   int
}

// We need to be able to accept multiple requests from the same connection
func (c *ConnHandler) Handle(ctx context.Context) {
	// handle error
	if c.destConn == nil {
		return
	} else if c.clientConn == nil {
		return
	}

	defer func() {
		c.destConn.Close()
		c.clientConn.Close()
	}()

	// Set deadlines for read and writes both ways
	c.destConn.SetDeadline(time.Now().Add(5000 * time.Millisecond))
	c.clientConn.SetDeadline(time.Now().Add(5000 * time.Millisecond))

	start := time.Now()

	raw, err := ReadFromConn(c.clientConn)
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

	n, err := io.Copy(c.destConn, bytes.NewReader(raw))
	if err != nil {
		c.Error(err)
		return
	}

	c.Log("[REQ-%d] forwarded request of %d bytes from %s to %s ", req.ID, n, httpReq.Host, c.destConn.RemoteAddr())

	// Take outgoing responses, mutate and forward
	body, err := ReadFromConn(c.destConn)
	if err != nil {
		c.Error(err)
		return
	}

	_, err = io.Copy(c.clientConn, bytes.NewReader(body))
	if err != nil {
		c.Fatal(err)
		return
	}

	c.Log("[REQ-%d] connection to %s closed (%dms)", req.ID, c.clientConn.RemoteAddr(), time.Now().Sub(start).Milliseconds())
}

func (c *ConnHandler) Log(msg string, args ...interface{}) {
	c.logger.Println(fmt.Sprintf("[CONN-%d]: %s", c.ID, fmt.Sprintf(msg, args...)))
}

func (c *ConnHandler) Error(err error) {
	c.Log(fmt.Sprintf("[ERR] %s", err.Error()))
}

func (c *ConnHandler) Fatal(err error) {
	c.Log(fmt.Sprintf("[ERR] %s", err.Error()))
}
