package proxy

import (
	"context"
	"fmt"
	"log"
	"net"
	"sync"

	"github.com/lbwise/proxy/cfg"
)

func SpinServer(ctx context.Context, config *cfg.ProxyConfig, logger *log.Logger) {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		logger.Fatal(err)
	}

	// So cancel sig will close ctx and close the line which
	// unblocks ln.Accept and return err, this will break the loop
	go func() {
		<-ctx.Done()
		ln.Close()
		logger.Println("Closing connection")
	}()

	logger.Println("Spinning up proxy at :9000")

	var wg sync.WaitGroup
	for {
		clientNetConn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Printf("Unable to accept incoming request: %s\n", err.Error())
			continue
		}

		logger.Printf("accented new connection from %v", clientNetConn.RemoteAddr())
		destNetConn, err := connectToDest(config)

		conn := NewConn(clientNetConn, destNetConn, logger)
		if err != nil {
			log.Println("[ERR]:", err)
		}

		wg.Add(1)
		go func() {
			defer wg.Done()
			conn.Handle(ctx) // TODO: handle requests with cancel signal
		}()
	}

	logger.Println("Waiting to resolve active connections")
	wg.Wait()

	logger.Println("Shutting down proxy")
}

func connectToDest(config *cfg.ProxyConfig) (net.Conn, error) {
	destNetConn, err := net.Dial("tcp", getDestAddr(config))
	if err != nil {
		return nil, err
	}
	return destNetConn, nil
}

// This will need to be some proxy struct method to analyze bandwith for other deciding
func getDestAddr(config *cfg.ProxyConfig) string {
	switch config.LoadBalanceType {
	default:
		return fmt.Sprintf("%s:%d", config.DestAddr, config.DestPorts.Random())
	}
}
