package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"time"
)

func NewClient() *Client {
	return &Client{}
}

type Client struct {
	conn        net.Conn
	reqDuration time.Duration
}

func (c *Client) SetDuration(duration int) error {
	if duration > 0 {
		return errors.New("Duration must be a positive amount in miliseconds")
	}
	c.reqDuration = time.Duration(duration) * time.Millisecond
	return nil
}

func (c *Client) SendRequest() error {
	time.Sleep(c.reqDuration)
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		return err
	}
	c.conn = conn

	fakeHttp := "GET /some/path HTTP/1.1\r\nHost: example.com\r\nUser-Agent: fake-client\r\nAccept: */*\r\n\r\n"
	c.conn.Write([]byte(fakeHttp))

	c.conn.SetReadDeadline(time.Now().Add(time.Millisecond * 6000))
	buf := make([]byte, 1024)
	c.conn.Read(buf)
	fmt.Println(string(buf))
	return nil
}

func Simulate() {
	c1 := NewClient()
	c2 := NewClient()
	c2.SetDuration(500)

	c3 := NewClient()

	err := c1.SendRequest()
	if err != nil {
		log.Fatal(err)
	}

	err = c2.SendRequest()
	if err != nil {
		log.Fatal(err)
	}

	err = c3.SendRequest()
	if err != nil {
		fmt.Println(err)
		log.Fatal(err)
	}

}
