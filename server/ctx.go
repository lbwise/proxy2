package server

import (
	"fmt"
	"log"
	"net"
	"strconv"
	"sync"
	"time"
)

func NewCtx(conn net.Conn, lg *log.Logger) *Ctx {
	return &Ctx{conn, lg}
}

type Ctx struct {
	conn net.Conn
	log  *log.Logger
}

func (c *Ctx) Log(msg string, args ...interface{}) {
	c.log.Println(fmt.Sprintf(msg, args...))
}

func (c *Ctx) Handle(wg *sync.WaitGroup) {
	c.Log("CONNECTED NEW CLIENT: %s", c.conn.RemoteAddr())
	time.Sleep(500 * time.Millisecond)

	port, _ := strconv.Atoi(c.conn.RemoteAddr().String()[10:])
	if port > 60100 {
		c.Log("INVALID PORT RECIEVED: %d", port)
		c.conn.Close()
		wg.Done()
		return
	}

	wg.Done()
	defer c.conn.Close()
}
