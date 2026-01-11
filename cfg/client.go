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
			{WaitBefore: 10 * time.Second, NumAgents: 100, ReqPath: "/ping"},
		},
	}

}

type ClientSimulationConfig struct {
	ClientAddr string          `yaml:"client-addr"`
	ProxyAddr  string          `yaml:"proxy-addr"`
	ProxyPort  Port            `yaml:"proxy-port"`
	PortRange  PortRange       `yaml:"port-range"`
	Flow       []CSInstruction `yaml:"flow"`
}

type CSInstruction struct {
	WaitBefore time.Duration `yaml:"wait-before"`
	NumAgents  int           `yaml:"num-agents"`
	ReqPath    string        `yaml:"req-path"`
	ReqBody    io.Reader     `yaml:"req-body"`
}
