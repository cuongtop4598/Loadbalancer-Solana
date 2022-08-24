package config

import (
	"fmt"
	"strings"

	"github.com/caarlos0/env/v6"
	"github.com/pkg/errors"
)

type Config struct {
	Port              int    `env:"PORT"`
	NodesString       string `env:"NODES"`
	Interval          int    `env:"CHECK_INTERVAL"`
	BlockThreshold    uint64 `env:"BLOCK_THRESHOLD"`
	ConnectionTimeout int    `env:"CONNECTION_TIMEOUT"`
	Nodes             []string
}

func ParseConfig() (Config, error) {
	c := Config{}
	if err := env.Parse(&c); err != nil {
		fmt.Printf("%+v\n", err)
		fmt.Println("Load config error!")
		return Config{}, err
	}
	c.Nodes = strings.Split(c.NodesString, ",")
	if len(c.Nodes) == 0 {
		return Config{}, errors.Errorf("Nodes are not defined")
	}
	return c, nil
}

func ParseConfigWPanic() Config {
	config, err := ParseConfig()

	if err != nil {
		panic(err)
	}

	return config
}
