package services

import (
	"encoding/hex"
	"strings"
)

var (
	transactionVersion                           = 1
	transactionProcessingTypeRelayed             = "RelayedTx"
	transactionProcessingTypeBuiltInFunctionCall = "BuiltInFunctionCall"
	transactionProcessingTypeMoveBalance         = "MoveBalance"
	builtInFunctionClaimDeveloperRewards         = "ClaimDeveloperRewards"
	refundGasMessage                             = "refundedGas"
	sendingValueToNonPayableContractDataPrefix   = "@" + hex.EncodeToString([]byte("sending value to non payable contract"))
	emptyHash                                    = strings.Repeat("0", 64)
)

var (
	transactionEventSignalError                 = "signalError"
	transactionEventESDTTransfer                = "ESDTTransfer"
	transactionEventESDTNFTTransfer             = "ESDTNFTTransfer"
	transactionEventMultiESDTNFTTransfer        = "MultiESDTNFTTransfer"
	transactionEventESDTLocalBurn               = "ESDTLocalBurn"
	transactionEventESDTLocalMint               = "ESDTLocalMint"
	transactionEventESDTWipe                    = "ESDTWipe"
	transactionEventTopicInvalidMetaTransaction = "meta transaction is invalid"
)
