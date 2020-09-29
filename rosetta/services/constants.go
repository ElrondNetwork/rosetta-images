package services

const (
	NumBlocksToGet = uint64(200)

	RosettaVersion = "1.4.0"

	// OpStatusOK is the operation status for successful operations.
	OpStatusSuccess = "Success"
	// OpStatusFailed is the operation status for failed operations.
	OpStatusFailed = "Failed"

	opTransfer = "Transfer"
	opFee      = "Fee"
	opReward   = "Reward"
	opScResult = "SmartContractResult"
)
