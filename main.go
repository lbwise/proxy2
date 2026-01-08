package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lbwise/proxy/client"
	"github.com/lbwise/proxy/dest"
	"github.com/lbwise/proxy/proxy"
)

/*
Okay so really we want in order:
set up dest srv -> pass config to proxy -> spin up proxy -> create clients -> simulate traffic
*/

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	baseLogger, closer, err := NewBaseLogger()
	if err != nil {
		log.Fatal(fmt.Sprintf("could not start proxy: %s", err.Error()))
		return
	}
	defer closer()

	wg.Add(3)

	// Spin up destination servers
	go func() {
		defer wg.Done()
		destLog := log.New(baseLogger.Writer(), "[DEST] ", baseLogger.Flags())
		_, err := dest.SpinServers(ctx, destLog)
		if err != nil {
			return
		}
	}()

	time.Sleep(time.Second)

	// Spin up proxy
	go func() {
		defer wg.Done()
		proxyLog := log.New(baseLogger.Writer(), "[PROXY] ", baseLogger.Flags())
		proxy.SpinServer(
			ctx,
			proxy.DefaultConfig(),
			proxyLog)
	}()

	time.Sleep(time.Second)

	// Simulate clients
	go func() {
		defer wg.Done()
		clientLog := log.New(baseLogger.Writer(), "[CLIENT] ", baseLogger.Flags())
		client.Simulate(ctx, clientLog)
	}()

	waitCloseSignal()
	cancel() // cancel context to let everything shut itself down
	wg.Wait()
}

func waitCloseSignal() {
	done := make(chan os.Signal, 1)
	signal.Notify(done, syscall.SIGINT, syscall.SIGTERM)
	<-done
}

func NewBaseLogger() (*log.Logger, func() error, error) {
	dtime := strings.Join(strings.Split(time.Now().String()[:19], " "), "-")
	f, err := os.Create(fmt.Sprintf("./logs/prox-server-%s", dtime))
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(f, "", log.LstdFlags)

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
