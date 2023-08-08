package testscommon

import (
	"fmt"
	"math/big"
	"strings"

	"github.com/multiversx/mx-chain-core-go/core"
	"github.com/multiversx/mx-chain-core-go/core/pubkeyConverter"
	"github.com/multiversx/mx-chain-core-go/data/api"
	"github.com/multiversx/mx-chain-core-go/data/transaction"
	"github.com/multiversx/mx-chain-go/sharding"
	"github.com/multiversx/mx-chain-proxy-go/data"
	"github.com/multiversx/mx-chain-rosetta/server/resources"
)

var (
	emptyHash        = strings.Repeat("0", 64)
	genesisTimestamp = int64(1596117600)
)

type networkProviderMock struct {
	pubKeyConverter core.PubkeyConverter

	MockIsOffline                   bool
	MockNumShards                   uint32
	MockObservedActualShard         uint32
	MockObservedProjectedShard      uint32
	MockObservedProjectedShardIsSet bool
	MockNativeCurrencySymbol        string
	MockCustomCurrencies            []resources.Currency
	MockGenesisBlockHash            string
	MockGenesisTimestamp            int64
	MockNetworkConfig               *resources.NetworkConfig
	MockGenesisBalances             []*resources.GenesisBalance
	MockNodeStatus                  *resources.AggregatedNodeStatus
	MockBlocksByNonce               map[uint64]*api.Block
	MockBlocksByHash                map[string]*api.Block
	MockNextAccountBlockCoordinates *resources.BlockCoordinates
	MockAccountsByAddress           map[string]*resources.Account
	MockAccountsNativeBalances      map[string]*resources.Account
	MockAccountsESDTBalances        map[string]*resources.AccountESDTBalance
	MockMempoolTransactionsByHash   map[string]*transaction.ApiTransactionResult
	MockComputedTransactionHash     string
	MockComputedReceiptHash         string
	MockNextError                   error

	SendTransactionCalled func(tx *data.Transaction) (string, error)
}

// NewNetworkProviderMock -
func NewNetworkProviderMock() *networkProviderMock {
	pubKeyConverter, _ := pubkeyConverter.NewBech32PubkeyConverter(32, log)

	return &networkProviderMock{
		pubKeyConverter:                 pubKeyConverter,
		MockIsOffline:                   false,
		MockNumShards:                   3,
		MockObservedActualShard:         0,
		MockObservedProjectedShard:      0,
		MockObservedProjectedShardIsSet: false,
		MockNativeCurrencySymbol:        "XeGLD",
		MockCustomCurrencies:            make([]resources.Currency, 0),
		MockGenesisBlockHash:            emptyHash,
		MockGenesisTimestamp:            genesisTimestamp,
		MockNetworkConfig: &resources.NetworkConfig{
			BlockchainName:         "MultiversX",
			NetworkID:              "T",
			NetworkName:            "testnet",
			MinGasPrice:            1000000000,
			MinGasLimit:            50000,
			GasPerDataByte:         1500,
			GasPriceModifier:       0.01,
			GasLimitCustomTransfer: 200000,
			ExtraGasLimitGuardedTx: 50000,
		},
		MockGenesisBalances: make([]*resources.GenesisBalance, 0),
		MockNodeStatus: &resources.AggregatedNodeStatus{
			Synced: true,
			LatestBlock: resources.BlockSummary{
				Nonce:             0,
				Hash:              emptyHash,
				PreviousBlockHash: emptyHash,
				Timestamp:         genesisTimestamp,
			},
			OldestBlockWithHistoricalState: resources.BlockSummary{
				Nonce:             0,
				Hash:              emptyHash,
				PreviousBlockHash: emptyHash,
				Timestamp:         genesisTimestamp,
			},
		},
		MockBlocksByNonce: make(map[uint64]*api.Block),
		MockBlocksByHash:  make(map[string]*api.Block),
		MockNextAccountBlockCoordinates: &resources.BlockCoordinates{
			Nonce:    0,
			Hash:     emptyHash,
			RootHash: emptyHash,
		},
		MockAccountsByAddress:         make(map[string]*resources.Account),
		MockAccountsNativeBalances:    make(map[string]*resources.Account),
		MockAccountsESDTBalances:      make(map[string]*resources.AccountESDTBalance),
		MockMempoolTransactionsByHash: make(map[string]*transaction.ApiTransactionResult),
		MockComputedTransactionHash:   emptyHash,
		MockNextError:                 nil,
	}
}

// IsOffline -
func (mock *networkProviderMock) IsOffline() bool {
	return mock.MockIsOffline
}

// GetBlockchainName -
func (mock *networkProviderMock) GetBlockchainName() string {
	return mock.MockNetworkConfig.BlockchainName
}

// GetNativeCurrency -
func (mock *networkProviderMock) GetNativeCurrency() resources.Currency {
	return resources.Currency{
		Symbol:   mock.MockNativeCurrencySymbol,
		Decimals: 18,
	}
}

// GetCustomCurrencies -
func (mock *networkProviderMock) GetCustomCurrencies() []resources.Currency {
	return mock.MockCustomCurrencies
}

// GetCustomCurrencyBySymbol -
func (mock *networkProviderMock) GetCustomCurrencyBySymbol(symbol string) (resources.Currency, bool) {
	for _, currency := range mock.MockCustomCurrencies {
		if currency.Symbol == symbol {
			return currency, true
		}
	}

	return resources.Currency{}, false
}

// HasCustomCurrency -
func (mock *networkProviderMock) HasCustomCurrency(symbol string) bool {
	_, has := mock.GetCustomCurrencyBySymbol(symbol)
	return has
}

// GetNetworkConfig -
func (mock *networkProviderMock) GetNetworkConfig() *resources.NetworkConfig {
	return mock.MockNetworkConfig
}

// GetGenesisBlockSummary -
func (mock *networkProviderMock) GetGenesisBlockSummary() *resources.BlockSummary {
	return &resources.BlockSummary{
		Nonce:     0,
		Hash:      mock.MockGenesisBlockHash,
		Timestamp: mock.MockGenesisTimestamp,
	}
}

// GetGenesisTimestamp -
func (mock *networkProviderMock) GetGenesisTimestamp() int64 {
	return mock.MockGenesisTimestamp
}

// GetGenesisBalances -
func (mock *networkProviderMock) GetGenesisBalances() ([]*resources.GenesisBalance, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	return mock.MockGenesisBalances, nil
}

// GetNodeStatus -
func (mock *networkProviderMock) GetNodeStatus() (*resources.AggregatedNodeStatus, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	return mock.MockNodeStatus, nil
}

// GetBlockByNonce -
func (mock *networkProviderMock) GetBlockByNonce(nonce uint64) (*api.Block, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	block, ok := mock.MockBlocksByNonce[nonce]
	if ok {
		return block, nil
	}

	return nil, fmt.Errorf("block %d not found", nonce)
}

// GetBlockByHash -
func (mock *networkProviderMock) GetBlockByHash(hash string) (*api.Block, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	block, ok := mock.MockBlocksByHash[hash]
	if ok {
		return block, nil
	}

	return nil, fmt.Errorf("block %s not found", hash)
}

// GetAccount -
func (mock *networkProviderMock) GetAccount(address string) (*resources.AccountOnBlock, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	account, ok := mock.MockAccountsByAddress[address]
	if ok {
		return &resources.AccountOnBlock{
			Account:          *account,
			BlockCoordinates: *mock.MockNextAccountBlockCoordinates,
		}, nil
	}

	return nil, fmt.Errorf("account %s not found", address)
}

func (mock *networkProviderMock) GetAccountNativeBalance(address string, _ resources.AccountQueryOptions) (*resources.AccountOnBlock, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	accountBalance, ok := mock.MockAccountsNativeBalances[address]
	if ok {
		return &resources.AccountOnBlock{
			Account: resources.Account{
				Balance: accountBalance.Balance,
				Nonce:   accountBalance.Nonce,
			},
			BlockCoordinates: *mock.MockNextAccountBlockCoordinates,
		}, nil
	}

	return nil, fmt.Errorf("account %s not found", address)
}

func (mock *networkProviderMock) GetAccountESDTBalance(address string, tokenIdentifier string, _ resources.AccountQueryOptions) (*resources.AccountESDTBalance, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	key := fmt.Sprintf("%s_%s", address, tokenIdentifier)
	accountBalance, ok := mock.MockAccountsESDTBalances[key]
	if ok {
		return &resources.AccountESDTBalance{
			Balance:          accountBalance.Balance,
			BlockCoordinates: *mock.MockNextAccountBlockCoordinates,
		}, nil
	}

	return nil, fmt.Errorf("account %s not found", address)
}

// IsAddressObserved -
func (mock *networkProviderMock) IsAddressObserved(address string) (bool, error) {
	if mock.MockNextError != nil {
		return false, mock.MockNextError
	}

	shardCoordinator, err := sharding.NewMultiShardCoordinator(mock.MockNumShards, mock.MockObservedActualShard)
	if err != nil {
		return false, err
	}

	pubKey, err := mock.ConvertAddressToPubKey(address)
	if err != nil {
		return false, err
	}

	shard := shardCoordinator.ComputeId(pubKey)

	isObservedActualShard := shard == mock.MockObservedActualShard
	isObservedProjectedShard := pubKey[len(pubKey)-1] == byte(mock.MockObservedProjectedShard)

	if mock.MockObservedProjectedShardIsSet {
		return isObservedProjectedShard, nil
	}

	return isObservedActualShard, nil
}

// ConvertPubKeyToAddress -
func (mock *networkProviderMock) ConvertPubKeyToAddress(pubkey []byte) string {
	return mock.pubKeyConverter.Encode(pubkey)
}

// ConvertAddressToPubKey -
func (mock *networkProviderMock) ConvertAddressToPubKey(address string) ([]byte, error) {
	return mock.pubKeyConverter.Decode(address)
}

// ComputeTransactionHash -
func (mock *networkProviderMock) ComputeTransactionHash(_ *data.Transaction) (string, error) {
	if mock.MockNextError != nil {
		return "", mock.MockNextError
	}

	return mock.MockComputedTransactionHash, nil
}

// ComputeReceiptHash -
func (mock *networkProviderMock) ComputeReceiptHash(_ *transaction.ApiReceipt) (string, error) {
	if mock.MockNextError != nil {
		return "", mock.MockNextError
	}

	return mock.MockComputedReceiptHash, nil
}

// ComputeTransactionFeeForMoveBalance -
func (mock *networkProviderMock) ComputeTransactionFeeForMoveBalance(tx *transaction.ApiTransactionResult) *big.Int {
	minGasLimit := mock.MockNetworkConfig.MinGasLimit
	extraGasLimitGuardedTx := mock.MockNetworkConfig.ExtraGasLimitGuardedTx
	gasPerDataByte := mock.MockNetworkConfig.GasPerDataByte
	gasLimit := minGasLimit + gasPerDataByte*uint64(len(tx.Data))

	isGuarded := len(tx.GuardianAddr) > 0
	if isGuarded {
		gasLimit += extraGasLimitGuardedTx
	}

	fee := core.SafeMul(gasLimit, tx.GasPrice)
	return fee
}

// SendTransaction -
func (mock *networkProviderMock) SendTransaction(tx *data.Transaction) (string, error) {
	if mock.MockNextError != nil {
		return "", mock.MockNextError
	}

	if mock.SendTransactionCalled != nil {
		return mock.SendTransactionCalled(tx)
	}

	return mock.MockComputedTransactionHash, nil
}

// GetMempoolTransactionByHash -
func (mock *networkProviderMock) GetMempoolTransactionByHash(hash string) (*transaction.ApiTransactionResult, error) {
	if mock.MockNextError != nil {
		return nil, mock.MockNextError
	}

	transactionObj, ok := mock.MockMempoolTransactionsByHash[hash]
	if ok {
		return transactionObj, nil
	}

	return nil, nil
}
