package core

import (
	"loadbalancer-solana/config"
	"time"

	"go.uber.org/zap"
)

func observeNode(node Node, config config.Config) Node {
	Log.Info("Observing node", zap.String("", node.Endpoint))
	blockNumber, err := getBlockNumber(&node, config)

	if err != nil {
		Log.Error("Obsserving failed with", zap.Any("", err))
		node.Available = false
	} else {
		node.Available = true
		node.BlockNumber = blockNumber
		Log.Info("Observing result", zap.Any("", node))
	}

	return node
}

func chooseBestNodeId(nodes []Node, config config.Config) (bestNodeId int) {
	var maxBlock uint64 = 0

	for i, n := range nodes {
		if n.Available && n.BlockNumber > maxBlock {
			maxBlock = n.BlockNumber

			bestNodeId = i
		}

	}

	return bestNodeId
}

func Observe(config config.Config, nodes []Node, currentNodeId *int) {
	for i, node := range nodes {
		nodes[i] = observeNode(node, config)
	}

	bestNodeId := chooseBestNodeId(nodes, config)

	if currentNodeId == nil {
		*currentNodeId = bestNodeId
	} else {
		currentNode := nodes[*currentNodeId]
		bestNode := nodes[bestNodeId]

		blockDiff := currentNode.BlockNumber - bestNode.BlockNumber

		if !currentNode.Available || blockDiff > config.BlockThreshold {
			*currentNodeId = bestNodeId
		}
	}

	Log.Info("Best node: ", zap.Any("", nodes[*currentNodeId].Endpoint))
}

func StartPeriodicObserve(config config.Config, nodes []Node, currentNodeId *int) {
	ticker := time.NewTicker(time.Duration(config.Interval) * time.Second)
	defer ticker.Stop()

	for range ticker.C {
		Observe(config, nodes, currentNodeId)
	}
}
