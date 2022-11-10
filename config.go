package main

import (
	"os"

	yaml "gopkg.in/yaml.v3"
)

type NodeConfig struct {
	TotalStorage    int64  `yaml:"total-storage"`
	ContractAddress string `yaml:"contract-address"`
	FactoryAddress  string `yaml:"factory-address"`
	TLSCert         string `yaml:"tls-cert"`
}

func LoadConfig() *NodeConfig {
	dat, err := os.ReadFile("./config.yaml")
	if err != nil {
		panic("Failed to open config.yaml")
	}
	config := NodeConfig{}
	err = yaml.Unmarshal(dat, &config)
	if err != nil {
		panic("Failed to unmarshal config")
	}

	return &config
}
