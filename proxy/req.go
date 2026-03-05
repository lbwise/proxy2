package proxy

import (
	"bytes"
	"net/http"
	"sync/atomic"
	"time"
)

func NewRequest(raw []byte) *Request {
	return &Request{
		ID:    GenNextReqId(),
		start: time.Now(),
		raw:   bytes.NewBuffer(raw),
		In:    &http.Request{},
	}
}

var (
	reqIdCounter uint64
)

func GenNextReqId() uint64 {
	return atomic.AddUint64(&reqIdCounter, 1)
}

type Request struct {
	ID       uint64
	Status   RequestStatus
	start    time.Time
	duration time.Duration
	raw      *bytes.Buffer
	In       *http.Request
	Out      *http.Request
}

type RequestStatus int

const (
	Active RequestStatus = iota
	Closed
	Error
)
