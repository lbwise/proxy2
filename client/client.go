package client

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net"
	"net/http"
	"sync"
	"time"
)

func NewClient(logger *log.Logger) *Client {
	return &Client{
		logger:   logger,
		deadline: 10 * time.Second,
	}
}

type Client struct {
	conn     net.Conn
	deadline time.Duration
	logger   *log.Logger
}

func (c *Client) SetDuration(duration int) error {
	if duration < 0 {
		return errors.New("duration must be a positive amount in miliseconds")
	}
	c.deadline = time.Duration(duration) * time.Millisecond
	return nil
}

func (c *Client) SendRequest(ctx context.Context) error {
	reqCtx, cancel := context.WithTimeout(ctx, c.deadline)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", "http://localhost:9000/ping", nil)
	if err != nil {
		return err
	}

	req.Host = "localhost:9000"

	reqErr := make(chan error, 1)
	go func() {
		c.logger.Println("Request being sent")
		res, err := http.DefaultClient.Do(req)

		if err != nil {
			reqErr <- err
			return
		}
		defer res.Body.Close()

		c.logger.Printf("Request completed with status: %d\n", res.StatusCode)
		reqErr <- nil
	}()

	select {
	case <-reqCtx.Done():
		return errors.New(fmt.Sprintf("request timed out - %s", reqCtx.Err().Error()))
	case err := <-reqErr:
		return err
	}
}

func Simulate(ctx context.Context, logger *log.Logger) {
	clients := make([]*Client, 100)
	for i := 0; i < 100; i++ {
		clients[i] = NewClient(logger)
	}

	var wg sync.WaitGroup
	for i, client := range clients {
		wg.Add(1)

		go func() {
			defer wg.Done()
			if err := client.SendRequest(ctx); err != nil {
				logger.Printf("Client %d failed with error: %s", i, err)
			}
		}()
	}

	wg.Wait()
	fmt.Println("SIMULATION COMPLETED")
}
