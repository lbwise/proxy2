package cfg

import (
	"fmt"
	"math/rand"
	"strconv"

	"gopkg.in/yaml.v3"
)

type Port int

func (p *Port) String() string {
	return strconv.Itoa(int(*p))
}

func (p *Port) UnmarshalYAML(value *yaml.Node) error {
	if value.Kind != yaml.ScalarNode {
		return fmt.Errorf("port must be a scalar")
	}

	i, err := strconv.Atoi(value.Value)
	if err != nil {
		return err
	}

	*p = Port(i)
	return nil
}

type PortRange struct {
	Start Port `yaml:"start"`
	End   Port `yaml:"end"`
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
	ProxyConfig            `yaml:"proxy"`
	DestConfig             `yaml:"dest"`
	ClientSimulationConfig `yaml:"client"`
	Logs                   []string `yaml:"logs"`
}
