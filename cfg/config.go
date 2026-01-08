package cfg

import (
	"math/rand"
	"strconv"
)

type Port int

func (p Port) String() string {
	return strconv.Itoa(int(p))
}

type PortRange struct {
	Start Port
	End   Port
}

func (r PortRange) Random() Port {
	return Port(rand.Intn(int(r.End)-int(r.Start)) + int(r.Start))
}

func (r PortRange) Next(cur Port) Port {
	if cur <= r.End {
		return cur + 1
	}
	return r.Start
}

type ProxySimConfig struct {
	*ClientSimulationConfig
	DestAddr      string
	DestPortRange PortRange
	ProxyConfig   *ProxyConfig
}
