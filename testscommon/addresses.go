package testscommon

import (
	"encoding/hex"
)

var (
	// TestAddressAlice is a test address
	TestAddressAlice = "erd1qyu5wthldzr8wx5c9ucg8kjagg0jfs53s8nr3zpz3hypefsdd8ssycr6th"
	// TestPubKeyAlice is a test pubkey
	TestPubKeyAlice, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressAlice)
	// TestPubKeyHexAlice is a test pubkey
	TestPubKeyHexAlice = hex.EncodeToString(TestPubKeyAlice)

	// TestAddressBob is a test address
	TestAddressBob = "erd1spyavw0956vq68xj8y4tenjpq2wd5a9p2c6j8gsz7ztyrnpxrruqzu66jx"
	// TestPubKeyBob is a test pubkey
	TestPubKeyBob, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressBob)
	// TestPubKeyHexBob is a test pubkey
	TestPubKeyHexBob = hex.EncodeToString(TestPubKeyBob)

	// TestAddressOfContract is a test address
	TestAddressOfContract = "erd1qqqqqqqqqqqqqpgqfejaxfh4ktp8mh8s77pl90dq0uzvh2vk396qlcwepw"
	// TestPubkeyOfContract is a test pubkey
	TestPubkeyOfContract, _ = RealWorldBech32PubkeyConverter.Decode(TestAddressOfContract)
	// TestPubkeyHexOfContract is a test pubkey
	TestPubkeyHexOfContract = hex.EncodeToString(TestPubkeyOfContract)
)
