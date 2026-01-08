package proxy

import (
	"bytes"
	"time"
)

func NewRequest(raw []byte) *Request {
	return &Request{
		ID:    1,
		start: time.Now(),
		raw:   bytes.NewBuffer(raw),
	}
}

type Request struct {
	ID       int64
	Status   RequestStatus
	start    time.Time
	duration time.Duration
	raw      *bytes.Buffer
}

type RequestStatus int

const (
	Active RequestStatus = iota
	Closed
	Error
)
