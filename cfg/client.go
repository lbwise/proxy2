package cfg

import (
	"io"
	"time"
)

func DefaultCSConfig() *ClientSimulationConfig {
	return &ClientSimulationConfig{
		ClientAddr: "localhost",
		ProxyAddr:  "localhost",
		ProxyPort:  8080,
		PortRange:  PortRange{5000, 5010},
		Flow: []CSInstruction{
			{WaitBefore: 0, NumAgents: 1, ReqPath: "/ping"},
			{WaitBefore: 2 * time.Second, NumAgents: 2, ReqPath: "/ping"},
			{WaitBefore: 0, NumAgents: 1, ReqPath: "/ping"},
			{WaitBefore: 5 * time.Second, NumAgents: 1, ReqPath: "/ping"},
		},
	}

}

type ClientSimulationConfig struct {
	ClientAddr string
	ProxyAddr  string
	ProxyPort  Port
	PortRange
	Flow []CSInstruction
}

type CSInstruction struct {
	WaitBefore time.Duration
	NumAgents  int
	ReqPath    string
	ReqBody    io.Reader
}
