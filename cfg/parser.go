package cfg

import (
	"os"

	"gopkg.in/yaml.v3"
)

func ParseCfgFile(filename string) (*ProxySimConfig, error) {
	fileBytes, err := os.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var cfg ProxySimConfig
	err = yaml.Unmarshal(fileBytes, &cfg)
	if err != nil {
		return nil, err
	}
	cfg.ProxyConfig.DestPorts = cfg.DestConfig.PortRange
	return &cfg, nil
}
