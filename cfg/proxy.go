package cfg

import (
	"fmt"
	"net"

	"gopkg.in/yaml.v3"
)

func DefaultProxyConfig(destAddr string, destPorts PortRange) *ProxyConfig {
	return &ProxyConfig{
		ProxyPort:       9000,
		LoadBalanceType: RoundRobin,
		IPBlacklist:     make([]net.IP, 0),
		DestAddr:        destAddr,
		DestPorts:       destPorts,
	}
}

type LoadBalanceType string

func (lb *LoadBalanceType) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("load balance type must be a scalar")
	}

	switch value.Value {
	case "round-robin":
		*lb = RoundRobin
	default:
		return fmt.Errorf("load balance type must be a round-robin")
	}
	return nil
}

const (
	RoundRobin LoadBalanceType = "RoundRobin"
	Random     LoadBalanceType = "Random"
	Weighted   LoadBalanceType = "Weighted"
)

type ProxyConfig struct {
	ProxyPort       int             `yaml:"port"`
	LoadBalanceType LoadBalanceType `yaml:"load-balance"`
	IPBlacklist     []net.IP        `yaml:"blacklist"`
	DestAddr        string          `yaml:"dest-addr"`
	DestPorts       PortRange       `yaml:"dest-ports"`
}
