package client

import (
	"errors"
	"log"
	"net"
	"net/http"
	"time"
)

func NewClient(logger *log.Logger) *Client {
	return &Client{logger: logger}
}

type Client struct {
	conn        net.Conn
	reqDuration time.Duration
	logger      *log.Logger
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

	req.Host = "localhost:9000"

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		c.logger.Println(err)
		return err
	}

	c.logger.Println("COMPLETED STATUS:", res.StatusCode)
	return nil
}

func Simulate(logger *log.Logger) {
	c1 := NewClient(logger)
	c2 := NewClient(logger)
	c3 := NewClient(logger)

	c2.SetDuration(500)

	err := c1.SendRequest()
	if err != nil {
		logger.Fatal(err)
	}

	err = c2.SendRequest()
	if err != nil {
		logger.Fatal(err)
	}

	err = c3.SendRequest()
	if err != nil {
		logger.Fatal(err)
	}

}
