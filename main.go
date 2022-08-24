package main

import (
	"loadbalancer-solana/config"
	"loadbalancer-solana/core"
	_ "loadbalancer-solana/core"

	"github.com/joho/godotenv"
	"go.uber.org/zap"
)

func initNodes(config config.Config) []core.Node {
	nodes := make([]core.Node, len(config.Nodes))

	for i, n := range config.Nodes {
		nodes[i] = core.Node{
			Endpoint:    n,
			BlockNumber: 0,
			Available:   false,
			RPCCounter:  0,
		}
	}
	return nodes
}

func main() {
	err := godotenv.Load()
	if err != nil {
		core.Log.Error("error loading .env file")
	}
	config := config.ParseConfigWPanic()
	core.Log.Info("Config", zap.Any("", config))
	nodes := initNodes(config)

	currentNodeId := 0

	core.Observe(config, nodes, &currentNodeId)
	go core.StartPeriodicObserve(config, nodes, &currentNodeId)

	core.StartProxy(config.Port, nodes, &currentNodeId)
}
