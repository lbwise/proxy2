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

	"github.com/lbwise/proxy/cfg"
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

// Simulate will have a list of instructions that each wait for each
func Simulate(ctx context.Context, config *cfg.ClientSimulationConfig, logger *log.Logger) {
	var wg sync.WaitGroup
	for i, instr := range config.Flow {
		logger.Printf("Running instruction: %d", i)
		wg.Add(1)

		time.Sleep(instr.WaitBefore)

		go func() {
			defer wg.Done()
			var instrWg sync.WaitGroup
			for i := 0; i < instr.NumAgents; i++ {
				instrWg.Add(1)
				go func() {
					defer instrWg.Done()
					client := NewClient(logger)
					if err := client.SendRequest(ctx, instr); err != nil {
						logger.Printf("Client %d failed with error: %s", i, err)
					}
				}()
			}
			instrWg.Wait()
		}()

		logger.Printf("Finishing instruction: %d", i)
	}

	wg.Wait()
	fmt.Println("SIMULATION COMPLETED")
}

func (c *Client) SendRequest(ctx context.Context, instr cfg.CSInstruction) error {
	reqCtx, cancel := context.WithTimeout(ctx, c.deadline)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, "GET", fmt.Sprintf("http://localhost:9000%s", instr.ReqPath), instr.ReqBody)
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
