package client

import (
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
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
	req, err := http.NewRequest("GET", "http://localhost:9000/ping", nil)
	if err != nil {
		return err
	}

	req.Header.Add("Host", "liamwise")

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err)
		return err
	}

	fmt.Printf("\nCOMPLETED STATUS: %s\n", res.StatusCode)
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
