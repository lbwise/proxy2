package server

import (
	"log"
	"net"
	"sync"
)

func NewDefaultProxyConfig(addr net.Addr) *ProxyConfig {
	return &ProxyConfig{addr}
}

type ProxyConfig struct {
	dest net.Addr
}

func SpinServer(config *ProxyConfig, wg *sync.WaitGroup, logger *log.Logger) {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("SPINNING UP SERVER AT :9000")

	destConn, err := net.Dial("tcp", config.dest.String())
	if err != nil {
		logger.Fatal(err)
	}
	logger.Println("CONNECTED TO DEST :%s", config.dest.String())

	wg.Add(1)
	go func() {
		for {
			srcConn, err := ln.Accept()
			ctx := NewCtx(srcConn, logger)
			if err != nil {
				log.Fatal(err)
			}

			wg.Add(1)
			go ctx.Handle(destConn, wg)
		}
	}()
}
