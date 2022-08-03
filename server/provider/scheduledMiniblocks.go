package provider

import (
	dataBlock "github.com/ElrondNetwork/elrond-go-core/data/block"
	"github.com/ElrondNetwork/elrond-proxy-go/data"
)

func (provider *networkProvider) simplifyBlockWithScheduledTransactions(block *data.Block) error {
	previousBlock, err := provider.doGetBlockByNonce(block.Nonce - 1)
	if err != nil {
		return err
	}

	nextBlock, err := provider.doGetBlockByNonce(block.Nonce + 1)
	if err != nil {
		return err
	}

	doSimplifyBlockWithScheduledTransactions(previousBlock, block, nextBlock)

	return nil
}

func hasOnlyNormalMiniblocks(block *data.Block) bool {
	for _, miniblock := range block.MiniBlocks {
		if miniblock.ProcessingType != dataBlock.Normal.String() {
			return false
		}
	}

	return true
}

func doSimplifyBlockWithScheduledTransactions(previousBlock *data.Block, block *data.Block, nextBlock *data.Block) {
	// Discard "processed" miniblocks in block N, since they already produced effects in N-1
	removeProcessedMiniblocksOfBlock(block)

	// Move "processed" miniblocks from N+1 to N
	processedMiniblocksInNextBlock := findProcessedMiniblocks(nextBlock)
	appendMiniblocksToBlock(block, processedMiniblocksInNextBlock)

	// Build an artificial miniblock holding the "invalid" transactions that produced their effects in block N,
	// and replace the existing (one or two "invalid" miniblocks).
	invalidTxs := gatherInvalidTransactions(previousBlock, block, nextBlock)
	invalidMiniblock := &data.MiniBlock{
		Type:         dataBlock.InvalidBlock.String(),
		Transactions: invalidTxs,
	}
	removeInvalidMiniblocks(block)

	if len(invalidMiniblock.Transactions) > 0 {
		appendMiniblocksToBlock(block, []*data.MiniBlock{invalidMiniblock})
	}

	// Discard "scheduled" miniblocks of N, since we've already brought the "processed" ones from N+1,
	// and also handled the "invalid" ones.
	removeScheduledMiniblocks(block)
}

func removeProcessedMiniblocksOfBlock(block *data.Block) {
	removeMiniblocksFromBlock(block, func(miniblock *data.MiniBlock) bool {
		return miniblock.ProcessingType == dataBlock.Processed.String()
	})
}

func removeScheduledMiniblocks(block *data.Block) {
	removeMiniblocksFromBlock(block, func(miniblock *data.MiniBlock) bool {
		hasProcessingTypeScheduled := miniblock.ProcessingType == dataBlock.Scheduled.String()
		hasConstructionStateNotFinal := miniblock.ConstructionState != dataBlock.Final.String()
		shouldRemove := hasProcessingTypeScheduled && hasConstructionStateNotFinal
		return shouldRemove
	})
}

func removeInvalidMiniblocks(block *data.Block) {
	removeMiniblocksFromBlock(block, func(miniblock *data.MiniBlock) bool {
		return miniblock.Type == dataBlock.InvalidBlock.String()
	})
}

func gatherInvalidTransactions(previousBlock *data.Block, block *data.Block, nextBlock *data.Block) []*data.FullTransaction {
	// Find "invalid" transactions that are "final" in N
	invalidTxsInBlock := findInvalidTransactions(block)
	// If also present in N-1, discard them
	scheduledTxsHashesPreviousBlock := findScheduledTransactionsHashes(previousBlock)
	invalidTxsInBlock = discardTransactions(invalidTxsInBlock, scheduledTxsHashesPreviousBlock)

	// Find "invalid" transactions in N+1 that are "scheduled" in N
	invalidTxsInNextBlock := findInvalidTransactions(nextBlock)
	scheduledTxsHashesInBlock := findScheduledTransactionsHashes(block)
	invalidTxsScheduledInBlock := filterTransactions(invalidTxsInNextBlock, scheduledTxsHashesInBlock)

	// Duplication might occur, since a block can contain two "invalid" miniblocks,
	// one added to block body, one saved in the receipts unit (at times, they have different content, different hashes).
	invalidTxs := append(invalidTxsInBlock, invalidTxsScheduledInBlock...)
	invalidTxs = deduplicateTransactions(invalidTxs)

	return invalidTxs
}

func findScheduledTransactionsHashes(block *data.Block) map[string]struct{} {
	txs := make(map[string]struct{})

	for _, miniblock := range block.MiniBlocks {
		hasProcessingTypeScheduled := miniblock.ProcessingType == dataBlock.Scheduled.String()
		hasConstructionStateNotFinal := miniblock.ConstructionState != dataBlock.Final.String()
		shouldAccumulateTxs := hasProcessingTypeScheduled && hasConstructionStateNotFinal

		if shouldAccumulateTxs {
			for _, tx := range miniblock.Transactions {
				txs[tx.Hash] = struct{}{}
			}
		}
	}

	return txs
}

func findProcessedMiniblocks(block *data.Block) []*data.MiniBlock {
	foundMiniblocks := make([]*data.MiniBlock, 0, len(block.MiniBlocks))

	for _, miniblock := range block.MiniBlocks {
		if miniblock.ProcessingType == dataBlock.Processed.String() {
			foundMiniblocks = append(foundMiniblocks, miniblock)
		}
	}

	return foundMiniblocks
}

func findInvalidTransactions(block *data.Block) []*data.FullTransaction {
	invalidTxs := make([]*data.FullTransaction, 0)

	for _, miniblock := range block.MiniBlocks {
		if miniblock.Type == dataBlock.InvalidBlock.String() {
			for _, tx := range miniblock.Transactions {
				invalidTxs = append(invalidTxs, tx)
			}
		}
	}

	return invalidTxs
}
