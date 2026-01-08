package cfg

import "net"

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

const (
	RoundRobin LoadBalanceType = "RoundRobin"
	Random     LoadBalanceType = "Random"
	Weighted   LoadBalanceType = "Weighted"
)

type ProxyConfig struct {
	ProxyPort int
	LoadBalanceType
	IPBlacklist []net.IP
	DestAddr    string
	DestPorts   PortRange
}
