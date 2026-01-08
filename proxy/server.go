package proxy

import (
	"context"
	"log"
	"net"
	"sync"
)

func DefaultConfig(addr net.Addr) *Config {
	return &Config{addr}
}

type Config struct {
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
*/

func SpinServer(ctx context.Context, config *Config, logger *log.Logger) {
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
		logger.Printf("accented new connection from %v", clientNetConn.RemoteAddr())
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			logger.Printf("Unable to accept incoming request: %s\n", err.Error())
			continue
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
		go func() {
			defer wg.Done()
			conn.Handle(ctx) // TODO: handle requests with cancel signal
		}()
	}

	logger.Println("Waiting to resolve active connections")
	wg.Wait()

	logger.Println("Shutting down proxy")
}
