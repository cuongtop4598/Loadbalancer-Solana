package core

import (
	"context"
	"loadbalancer-solana/config"

	"github.com/gagliardetto/solana-go/rpc"
)

func getBlockNumber(node *Node, config config.Config) (uint64, error) {
	client := rpc.New(node.Endpoint)
	blockLastest, err := client.GetRecentBlockhash(
		context.TODO(),
		rpc.CommitmentFinalized,
	)
	if err != nil {
		Log.Error(err.Error())
	}

	endSlot := uint64(blockLastest.Context.Slot)
	blockNumber, err := client.GetBlocks(
		context.TODO(),
		uint64(blockLastest.Context.Slot),
		&endSlot,
		rpc.CommitmentFinalized,
	)
	if err != nil {
		Log.Error(err.Error())
	}
	if len(blockNumber) > 0 {
		return blockNumber[0], nil
	}
	return 0, nil
}
