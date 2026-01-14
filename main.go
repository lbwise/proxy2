package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/lbwise/proxy/cfg"
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

	cf, err := cfg.ParseCfgFile("./cfg/test.yaml")
	if err != nil {
		baseLogger.Fatal(fmt.Sprintf("could not parse config: %s", err.Error()))
		return
	}

	// Spin up destination servers
	go func() {
		defer wg.Done()
		destLog := log.New(baseLogger.Writer(), "[DEST] ", baseLogger.Flags())
		if !cf.LogContains("dest") {
			destLog.SetOutput(io.Discard)
		}
		_, err := dest.SpinServers(ctx, cf.DestConfig, destLog)
		if err != nil {
			return
		}
	}()

	time.Sleep(time.Second)

	// Spin up proxy
	go func() {
		defer wg.Done()

		proxyLog := log.New(baseLogger.Writer(), "[PROXY] ", baseLogger.Flags())
		if !cf.LogContains("proxy") {
			proxyLog.SetOutput(io.Discard)
		}
		p := proxy.New(&cf.ProxyConfig, proxyLog)
		p.SpinServer(ctx)
	}()

	time.Sleep(time.Second)

	// Simulate clients
	go func() {
		defer wg.Done()
		clientLog := log.New(baseLogger.Writer(), "[CLIENT] ", baseLogger.Flags())
		if !cf.LogContains("client") {
			clientLog.SetOutput(io.Discard)
		}
		client.Simulate(ctx, cf.ClientSimulationConfig, clientLog)
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
	dTime := strings.Join(strings.Split(time.Now().String()[:19], " "), "-")
	f, err := os.Create(fmt.Sprintf("./logs/prox-server-%s", dTime))
	if err != nil {
		return nil, nil, err
	}

	logger := log.New(f, "", log.LstdFlags)

	return logger, f.Close, nil
}
