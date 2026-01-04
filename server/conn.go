package server

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"sync"
	"time"

	"github.com/bytedance/gopkg/util/logger"
	_ "github.com/lbwise/learning/compressor/algos"
	"github.com/ugorji/go/codec"
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

func (c *Conn) Handle(wg *sync.WaitGroup) {
	// handle error
	if c.DestConn != nil {
		return
	} else if c.ClientConn != nil {
		return
	}

	defer func() {
		c.DestConn.Close()
		c.ClientConn.Close()
		wg.Done()
	}()

	start := time.Now()

	raw, err := ReadFromConn(c.ClientConn)
	if err != nil {
		c.Error(err)
		return
	}

	req := NewRequest(raw)

	// Take incoming requests, mutate and forward
	req, err :=

		c.Log("[REQ] %s %s host=%s agent=%s", req.Method, req.URL, req.Host, req.UserAgent())
	if err != nil {
		return
	}
	c.req = req

	buf := new(bytes.Buffer)
	err = c.req.Write(buf)
	if err != nil {
		c.Fatal(err)
		return
	}

	err = c.forward(buf, destConn)
	if err != nil {
		c.Fatal(err)
		return
	}

	// Take outgoing responses, mutate and forward
	destConn.SetDeadline(time.Now().Add(5000 * time.Millisecond))
	c.res, err = c.ingestRes(destConn)
	if err != nil {
		c.Fatal(err)
		return
	}

	resBuf := new(bytes.Buffer)
	err = c.res.Write(resBuf)
	if err != nil {
		c.Fatal(err)
		return
	}

	err = c.forward(resBuf, c.conn)
	if err != nil {
		c.Fatal(err)
		return
	}

	c.Log("Connection to %s closed (%dms)", c.conn.RemoteAddr(), time.Now().Sub(start).Milliseconds())
}

func BufferToHttpReq(r io.Reader) (*http.Request, error) {
	req, err := http.ReadRequest(bytes.Buffer{})
	if err != nil {
		return nil, err
	}

	return req, nil
}

func (c *Conn) ingestRes(input io.Reader) (*http.Response, error) {
	r := bufio.NewReader(input)
	res, err := http.ReadResponse(r, c.req)
	if err != nil {
		return nil, err
	}

	//c.Log("[RES] %s %s host=%s agent=%s", res.Method, res.URL, req.Host, req.UserAgent())

	return res, nil
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
