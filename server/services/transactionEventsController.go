package services

import (
	"strings"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
)

type transactionEventsController struct {
	provider NetworkProvider
}

func newTransactionEventsController(provider NetworkProvider) *transactionEventsController {
	return &transactionEventsController{
		provider: provider,
	}
}

func (controller *transactionEventsController) hasSignalErrorOfSendingValueToNonPayableContract(tx *transaction.ApiTransactionResult) bool {
	if !controller.hasEvents(tx) {
		return false
	}

	for _, event := range tx.Logs.Events {
		isSignalError := event.Identifier == transactionEventSignalError
		dataAsString := string(event.Data)
		dataMatchesError := strings.HasPrefix(dataAsString, sendingValueToNonPayableContractDataPrefix)

		if isSignalError && dataMatchesError {
			return true
		}
	}

	return false
}

func (controller *transactionEventsController) hasSignalErrorOfMetaTransactionIsInvalid(tx *transaction.ApiTransactionResult) bool {
	if !controller.hasEvents(tx) {
		return false
	}

	for _, event := range tx.Logs.Events {
		isSignalError := event.Identifier == transactionEventSignalError
		if !isSignalError {
			continue
		}

		if eventHasTopic(event, transactionEventTopicInvalidMetaTransaction) {
			return true
		}
		if eventHasTopic(event, transactionEventTopicInvalidMetaTransactionNotEnoughGas) {
			return true
		}
	}

	return false
}

func (controller *transactionEventsController) hasEvents(tx *transaction.ApiTransactionResult) bool {
	return tx.Logs != nil && tx.Logs.Events != nil && len(tx.Logs.Events) > 0
}

func eventHasTopic(event *transaction.Events, topicToFind string) bool {
	for _, topic := range event.Topics {
		if string(topic) == topicToFind {
			return true
		}
	}

	return false
}
