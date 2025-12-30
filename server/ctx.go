package server

import (
	"errors"
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"

	_ "github.com/lbwise/learning/compressor/algos"
	"github.com/lbwise/proxy/http"
)

func NewCtx(conn net.Conn, lg *log.Logger) *Ctx {
	return &Ctx{conn, lg, 1024}
}

type Ctx struct {
	conn     net.Conn
	log      *log.Logger
	readSize int
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

func (c *Ctx) Handle(wg *sync.WaitGroup) {
	c.Log("CONNECTED NEW CLIENT: %s", c.conn.RemoteAddr())
	time.Sleep(500 * time.Millisecond)
	err := c.WriteString("HELLO BACK")
	if err != nil {
		c.Error(err)
	}

	msg, err := c.Read()
	if err != nil {
		c.Fatal(err)
		return
	}
	parser := http.NewParser(msg)
	httpReq, err := parser.Parse()
	if err != nil {
		c.Fatal(err)
		return
	}

	c.Log("REQUEST METHOD: %s", httpReq.String())

	//compMsg, err := comp.NewCompressionAlgo(comp.RLEType, comp.NewAlgoConfig()).Compress(msg)
	//if err != nil {
	//	c.Fatal(err)
	//	return
	//}

	c.Log("COMP MSG RECEIVED: %s", msg)

	port, _ := strconv.Atoi(c.conn.RemoteAddr().String()[10:])
	if port > 60100 {
		c.Fatal(errors.New(fmt.Sprintf("INVALID PORT RECEIVED: %d", port)))
		wg.Done()
		return
	}

	wg.Done()
	defer c.conn.Close()
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
