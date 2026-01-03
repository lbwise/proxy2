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

	_ "github.com/lbwise/learning/compressor/algos"
)

func NewCtx(conn net.Conn, lg *log.Logger) *Ctx {
	return &Ctx{conn, lg, 1024, nil, nil}
}

type Ctx struct {
	conn     net.Conn
	log      *log.Logger
	readSize int
	req      *http.Request
	res      *http.Response
}

func (c *Ctx) Write(buf []byte) error {
	_, err := c.conn.Write(buf)
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctx) WriteString(msg string) error {
	_, err := c.conn.Write([]byte(msg))
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctx) Read() ([]byte, error) {
	buf := make([]byte, c.readSize)
	n, err := c.conn.Read(buf)
	if err != nil {
		return nil, err
	}
	return buf[:n], nil
}

func (c *Ctx) ReadString() (string, error) {
	buf, err := c.Read()
	return string(buf), err
}

func (c *Ctx) Handle(destConn net.Conn, wg *sync.WaitGroup) {
	defer func() {
		c.conn.Close()
		wg.Done()
	}()

	// Take incoming requests, mutate and forward
	input, err := c.Read()
	if err != nil {
		c.Fatal(err)
		return
	}

	req, err := c.ingestReq(input)
	if err != nil {
		return
	}
	c.req = req

	buf := new(bytes.Buffer)
	err = req.Write(buf)
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

	res, err := c.ingestRes(destConn)
	if err != nil {
		c.Fatal(err)
		return
	}

	resBuf := new(bytes.Buffer)
	err = res.Write(resBuf)
	if err != nil {
		c.Fatal(err)
		return
	}

	err = c.forward(resBuf, c.conn)
	if err != nil {
		c.Fatal(err)
		return
	}
}

func (c *Ctx) ingestReq(input []byte) (*http.Request, error) {
	r := bytes.NewBuffer(input)
	req, err := http.ReadRequest(bufio.NewReader(r))
	if err != nil {
		c.Fatal(err)
		return nil, err
	}

	c.Log("REQEUST METHOD: ", req.Method)

	return req, nil
}

func (c *Ctx) ingestRes(input io.Reader) (*http.Response, error) {
	r := bufio.NewReader(input)
	res, err := http.ReadResponse(r, c.req)
	if err != nil {
		fmt.Println("THIS ERROR HERE", err)
		return nil, err
	}
	return res, nil
}

func (c *Ctx) forward(r io.Reader, w io.Writer) error {
	_, err := io.Copy(w, r)
	if err != nil {
		return err
	}
	return nil
}

func (c *Ctx) Log(msg string, args ...interface{}) {
	c.log.Println(fmt.Sprintf(msg, args...))
}

func (c *Ctx) Error(err error) {
	c.log.Println(err.Error())
}

func (c *Ctx) Fatal(err error) {
	c.log.Println(err.Error())
	c.conn.Close()
}

/*
so what is the proxy actually doing

on setup its:
1. starting server
2. connecting to dest
3. listening for conns and handling them

on each request it:
1. reads from the conn
2. analyzes the req
3. forwards the request to dest
4. waits for a response
5. analyzes the response
6. forwards res to client
*/
