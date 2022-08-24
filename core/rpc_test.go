package core

import (
	"context"
	"fmt"
	"testing"

	"github.com/gagliardetto/solana-go/rpc"
)

func TestRPC(t *testing.T) {
	client := rpc.New("http://localhost:8000")
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
	fmt.Println(blockNumber)
}
