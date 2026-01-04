package server

import (
	"log"
	"net"
	"sync"
	"sync/atomic"
)

func NewDefaultProxyConfig(addr net.Addr) *ProxyConfig {
	return &ProxyConfig{addr}
}

type ProxyConfig struct {
	dest net.Addr
}

/*
TODO:
[] properly implement conn id and requeset id
[] add header injection and passing blocking
[] add analytics for proxy
[] add multiple destination server spin
[] add active connections and rate limiting


proxy will now need to manage lots of things like (ideally)
request id
conn id - x
active connections count -> rate limiting
header forwarding
blocking admin routes
req latency
load balancing
config file for proxy setup and handling multiple dest servers/clients/connections

conn obj:
- conn-id
- net conn obj
- []req obj
	- id
	- latency
	- status
	- conn-id
	- start
- closed
- last
*/

func SpinServer(config *ProxyConfig, wg *sync.WaitGroup, logger *log.Logger) {
	var counter atomic.Int64
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("SPINNING UP SERVER AT :9000")

	logger.Println("CONNECTED TO DEST: ", config.dest.String())

	wg.Add(1)
	go func() {
		for {
			clientNetConn, err := ln.Accept()
			if err != nil {
				logger.Println(err)
			}

			destNetConn, err := net.Dial("tcp", config.dest.String())
			if err != nil {
				logger.Println(err)
			}

			conn := NewConn(clientNetConn, destNetConn, logger)
			if err != nil {
				log.Fatal(err)
			}

			wg.Add(1)
			go conn.Handle(wg)
		}
	}()
}
