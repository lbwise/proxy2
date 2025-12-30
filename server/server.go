package server

import (
	"log"
	"net"
	"sync"
)

func SpinServer(wg *sync.WaitGroup, logger *log.Logger) {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatal(err)
	}

	logger.Println("SPINNING UP SERVER")

	wg.Add(1)
	go func() {
		for {
			conn, err := ln.Accept()
			ctx := NewCtx(conn, logger)
			if err != nil {
				log.Fatal(err)
			}

			wg.Add(1)
			go ctx.Handle(wg)
		}
	}()
}
