package services

import (
	"testing"

	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-rosetta/testscommon"
	"github.com/stretchr/testify/require"
)

func TestTransactionEventsController_HasAnySignalError(t *testing.T) {
	networkProvider := testscommon.NewNetworkProviderMock()
	controller := newTransactionEventsController(networkProvider)

	t.Run("arbitrary tx", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{}
		txMatches := controller.hasAnySignalError(tx)
		require.False(t, txMatches)
	})

	t.Run("tx with event 'signalError'", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{

				Events: []*transaction.Events{
					{
						Identifier: transactionEventSignalError,
					},
				},
			},
		}

		txMatches := controller.hasAnySignalError(tx)
		require.True(t, txMatches)
	})
}

func TestTransactionEventsController_FindManyEventsByIdentifier(t *testing.T) {
	networkProvider := testscommon.NewNetworkProviderMock()
	controller := newTransactionEventsController(networkProvider)

	t.Run("no matching events", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Events: []*transaction.Events{
					{
						Identifier: "a",
					},
				},
			},
		}

		events := controller.findManyEventsByIdentifier(tx, "b")
		require.Len(t, events, 0)
	})

	t.Run("more than one matching event", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Events: []*transaction.Events{
					{
						Identifier: "a",
						Data:       []byte("1"),
					},
					{
						Identifier: "a",
						Data:       []byte("2"),
					},
					{
						Identifier: "b",
						Data:       []byte("3"),
					},
				},
			},
		}

		events := controller.findManyEventsByIdentifier(tx, "a")
		require.Len(t, events, 2)
		require.Equal(t, []byte("1"), events[0].Data)
		require.Equal(t, []byte("2"), events[1].Data)
	})
}

func TestTransactionEventsController_HasSignalErrorOfSendingValueToNonPayableContract(t *testing.T) {
	networkProvider := testscommon.NewNetworkProviderMock()
	controller := newTransactionEventsController(networkProvider)

	t.Run("arbitrary tx", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{}
		txMatches := controller.hasSignalErrorOfSendingValueToNonPayableContract(tx)
		require.False(t, txMatches)
	})

	t.Run("invalid tx with event 'sending value to non-payable contract'", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{

				Events: []*transaction.Events{
					{
						Identifier: transactionEventSignalError,
						Data:       []byte(sendingValueToNonPayableContractDataPrefix + "aaaabbbbccccdddd"),
					},
				},
			},
		}

		txMatches := controller.hasSignalErrorOfSendingValueToNonPayableContract(tx)
		require.True(t, txMatches)
	})
}

func TestTransactionEventsController_HasSignalErrorOfMetaTransactionIsInvalid(t *testing.T) {
	networkProvider := testscommon.NewNetworkProviderMock()
	controller := newTransactionEventsController(networkProvider)

	t.Run("arbitrary tx", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{}
		txMatches := controller.hasSignalErrorOfMetaTransactionIsInvalid(tx)
		require.False(t, txMatches)
	})

	t.Run("invalid tx with event 'meta transaction is invalid'", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{

				Events: []*transaction.Events{
					{
						Identifier: transactionEventSignalError,
						Topics: [][]byte{
							[]byte(transactionEventTopicInvalidMetaTransaction),
						},
					},
				},
			},
		}

		txMatches := controller.hasSignalErrorOfMetaTransactionIsInvalid(tx)
		require.True(t, txMatches)
	})

	t.Run("invalid tx with event 'meta transaction is invalid: not enough gas'", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{

				Events: []*transaction.Events{
					{
						Identifier: transactionEventSignalError,
						Topics: [][]byte{
							[]byte(transactionEventTopicInvalidMetaTransactionNotEnoughGas),
						},
					},
				},
			},
		}

		txMatches := controller.hasSignalErrorOfMetaTransactionIsInvalid(tx)
		require.True(t, txMatches)
	})
}

func TestEventHasTopic(t *testing.T) {
	event := transaction.Events{
		Identifier: transactionEventSignalError,
		Topics: [][]byte{
			[]byte("foo"),
		},
	}

	require.True(t, eventHasTopic(&event, "foo"))
	require.False(t, eventHasTopic(&event, "bar"))
}

func TestTransactionEventsController(t *testing.T) {
	networkProvider := testscommon.NewNetworkProviderMock()
	controller := newTransactionEventsController(networkProvider)

	t.Run("ESDTNFTCreate", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Events: []*transaction.Events{
					{
						Identifier: "ESDTNFTCreate",
						Address:    testscommon.TestAddressAlice,
						Topics: [][]byte{
							[]byte("EXAMPLE-abcdef"),
							{0x2a},
							{0x1},
							{0x0},
						},
					},
				},
			},
		}

		events, err := controller.extractEventsESDTNFTCreate(tx)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "EXAMPLE-abcdef", events[0].identifier)
		require.Equal(t, testscommon.TestAddressAlice, events[0].otherAddress)
		require.Equal(t, []byte{0x2a}, events[0].nonceAsBytes)
		require.Equal(t, "1", events[0].value)
	})

	t.Run("ESDTNFTBurn", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Events: []*transaction.Events{
					{
						Identifier: "ESDTNFTBurn",
						Address:    testscommon.TestAddressAlice,
						Topics: [][]byte{
							[]byte("EXAMPLE-abcdef"),
							{0x2a},
							{0x1},
						},
					},
				},
			},
		}

		events, err := controller.extractEventsESDTNFTBurn(tx)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "EXAMPLE-abcdef", events[0].identifier)
		require.Equal(t, testscommon.TestAddressAlice, events[0].otherAddress)
		require.Equal(t, []byte{0x2a}, events[0].nonceAsBytes)
		require.Equal(t, "1", events[0].value)
	})

	t.Run("ESDTNFTAddQuantity", func(t *testing.T) {
		tx := &transaction.ApiTransactionResult{
			Logs: &transaction.ApiLogs{
				Events: []*transaction.Events{
					{
						Identifier: "ESDTNFTAddQuantity",
						Address:    testscommon.TestAddressAlice,
						Topics: [][]byte{
							[]byte("EXAMPLE-aabbcc"),
							{0x2a},
							{0x64},
						},
					},
				},
			},
		}

		events, err := controller.extractEventsESDTNFTAddQuantity(tx)
		require.NoError(t, err)
		require.Len(t, events, 1)
		require.Equal(t, "EXAMPLE-aabbcc", events[0].identifier)
		require.Equal(t, testscommon.TestAddressAlice, events[0].otherAddress)
		require.Equal(t, []byte{0x2a}, events[0].nonceAsBytes)
		require.Equal(t, "100", events[0].value)
	})
}
