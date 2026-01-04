package main

import (
	"fmt"
	"log"
	"net"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/lbwise/proxy/client"
	"github.com/lbwise/proxy/dest"
	"github.com/lbwise/proxy/server"
)

/*
Okay so really we want in order:
set up dest srv -> pass config to proxy -> spin up proxy -> create clients -> simulate traffic
*/

func main() {
	baseLogger, closer, err := NewBaseLogger()
	if err != nil {
		log.Fatal(err)
		return
	}
	defer closer()

	wg := &sync.WaitGroup{}
	wg.Add(1)

	destLog := log.New(baseLogger.Writer(), "[DEST] ", baseLogger.Flags())
	destServers, err := dest.SpinServers(destLog)

	proxyLog := log.New(baseLogger.Writer(), "[PROXY] ", baseLogger.Flags())
	server.SpinServer(
		server.NewDefaultProxyConfig(&net.TCPAddr{
			IP:   net.ParseIP(destServers[0].Addr),
			Port: destServers[0].Port,
		}),
		wg, proxyLog)

	clientLog := log.New(baseLogger.Writer(), "[CLIENT] ", baseLogger.Flags())
	client.Simulate(clientLog)
	wg.Done()

	wg.Wait()
	proxyLog.Println("GRACEFULLY CLOSING SERVER")
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
