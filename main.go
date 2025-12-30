package main

import (
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lbwise/proxy/client"
	"github.com/lbwise/proxy/server"
)

func main() {
	logger, closer, err := NewLogger()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer closer()

	wg := &sync.WaitGroup{}
	server.SpinServer(wg, logger)

	client.Simulate()
	wg.Done()

	wg.Wait()
	logger.Println("GRACEFULLY CLOSING SERVER")
}

func NewLogger() (*log.Logger, func() error, error) {
	dtime := strings.Join(strings.Split(time.Now().String()[:19], " "), "-")
	f, err := os.Create(fmt.Sprintf("./logs/prox-server-%s", dtime))
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(f, "SERVER: ", log.LstdFlags)

	return logger, f.Close, nil
}

/*
what does a proxy do?

takes incoming requests, and at each layer anaylzes and forward/blocks requests
configurable from the user (config file?)
logs traffic and latency
content based routing

examples:
- HTTP layer: inspects content
- TCP layer: blocks ports
*/
