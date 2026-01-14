package proxy

import (
	"context"
	"log"
	"net"
	"sync"

	"github.com/lbwise/proxy/cfg"
)

func New(config *cfg.ProxyConfig, logger *log.Logger) *Proxy {
	//for _, port := range config.DestPorts {
	//
	//}
	//
	return &Proxy{
		logger:        logger,
		config:        config,
		destConnPools: make(map[string]*DestConnectionPool),
	}
}

func (p *Proxy) DialDestServers() error {
	for _, dest := range p.config.DestPorts.ToArray() {
		addr := &net.TCPAddr{IP: net.IP(p.config.DestAddr), Port: dest}
		pool, err := NewDestConnectionPool(addr, 10)
		if err != nil {
			// Should handle better and not be fatal
			return err
		}
		p.destConnPools[addr.String()] = pool
	}
	return nil
}

type Proxy struct {
	logger        *log.Logger
	config        *cfg.ProxyConfig
	destConnPools map[string]*DestConnectionPool
}

func (p *Proxy) SpinServer(ctx context.Context) {
	ln, err := net.Listen("tcp", ":9000")
	if err != nil {
		p.logger.Fatal(err)
	}

	// So cancel sig will close ctx and close the line which
	// unblocks ln.Accept and return err, this will break the loop
	go func() {
		<-ctx.Done()
		ln.Close()
		p.logger.Println("Closing connection")
	}()

	p.logger.Println("Spinning up proxy at :9000")

	err = p.DialDestServers()
	if err != nil {
		p.logger.Fatal(err)
		return
	}
	p.logger.Println("CONNECTION POOL CREATED")
	p.logger.Println(p.destConnPools)

	var wg sync.WaitGroup
	for {
		clientNetConn, err := ln.Accept()
		if err != nil {
			if ctx.Err() != nil {
				break
			}
			p.logger.Printf("Unable to accept incoming request: %s\n", err.Error())
			continue
		}

		if p.isBlacklistedIP(clientNetConn) {
			p.logger.Println("Blacklisted client received, rejected connection")
			clientNetConn.Close()
			return
		}

		p.logger.Printf("accented new connection from %v", clientNetConn.RemoteAddr())
		destConnAddr, destConn := p.connectToDest()

		conn := NewConn(clientNetConn, destConn, p.logger)
		if err != nil {
			p.logger.Println("[ERR]:", err)
		}

		wg.Add(1)
		go func() {
			defer p.destConnPools[destConnAddr.String()].Release(destConn) // This is so gross make this nicer
			defer wg.Done()
			conn.Handle(ctx) // TODO: handle requests with cancel signal
		}()

	}

	p.logger.Println("Waiting to resolve active connections")
	wg.Wait()

	p.logger.Println("Shutting down proxy")
}

// Fetch connection from the conn pool of the chosen server
// The chosen server is done by load balancing method
// function will block until connection is available
func (p *Proxy) connectToDest() (net.Addr, net.Conn) {
	addr := p.getDestAddr()
	pool := p.destConnPools[addr.String()]
	return addr, pool.Acquire()
}

// This will need to be some proxy struct method to analyze bandwidth for other deciding
func (p *Proxy) getDestAddr() net.Addr {
	switch p.config.LoadBalanceType {
	default:
		return &net.TCPAddr{
			IP:   net.IP(p.config.DestAddr),
			Port: int(p.config.DestPorts.Random()),
		}
	}
}

func (p *Proxy) isBlacklistedIP(connect net.Conn) bool {
	for _, blIp := range p.config.IPBlacklist {
		if connect.RemoteAddr().String() == blIp.String() {
			return true
		}
	}
	return false
}
