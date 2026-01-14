package proxy

import (
	"net"
)

func NewDestConnectionPool(addr net.Addr, size uint8) (*DestConnectionPool, error) {
	connCh := make(chan net.Conn, size)

	for i := 0; i < int(size); i++ {
		conn, err := net.Dial("tcp", addr.String())
		if err != nil {
			return nil, err
		}
		connCh <- conn
	}

	return &DestConnectionPool{
		Addr:  addr,
		size:  size,
		conns: connCh,
	}, nil
}

type DestConnectionPool struct {
	Addr  net.Addr
	conns chan net.Conn
	size  uint8
}

func (p *DestConnectionPool) Acquire() net.Conn {
	return <-p.conns
}

func (p *DestConnectionPool) Release(conn net.Conn) {
	p.conns <- conn
}
