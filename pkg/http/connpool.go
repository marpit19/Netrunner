package http

import (
	"net"
	"sync"
)

type ConnPool struct {
	mu       sync.Mutex
	conns    chan net.Conn
	maxConns int
}

func NewConnPool(maxConns int) *ConnPool {
	return &ConnPool{
		conns:    make(chan net.Conn, maxConns),
		maxConns: maxConns,
	}
}

func (p *ConnPool) Get(network, address string) (net.Conn, error) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case conn := <-p.conns:
		return conn, nil
	default:
		return net.Dial(network, address)
	}
}

func (p *ConnPool) Put(conn net.Conn) {
	p.mu.Lock()
	defer p.mu.Unlock()

	select {
	case p.conns <- conn:
	default:
		conn.Close()
	}
}

func (p *ConnPool) CloseIdleConnections() {
	p.mu.Lock()
	defer p.mu.Unlock()

	close(p.conns)
	for conn := range p.conns {
		conn.Close()
	}
	p.conns = make(chan net.Conn, p.maxConns)
}
