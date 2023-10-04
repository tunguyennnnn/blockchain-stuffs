package wallets

import (
	"encoding/hex"
	"log"
	"os"

	"github.com/blocto/solana-go-sdk/types"
	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type SeedPhraseSize int

type WalletChildKeys struct {
	PrivateKey    string
	XCoordinate   string
	YCoordinate   string
	PublicAddress string
}

type BitcoinChildKeys struct {
	PrivateKey    string
	XCoordinate   string
	YCoordinate   string
	P2pkhAddress  string
	SegwitAddress string
}

const (
	SeedPhraseSize12 SeedPhraseSize = 12
	SeedPhraseSize24 SeedPhraseSize = 24
)

type WalletKey struct {
	SeedPhraseSize SeedPhraseSize
	SeedPhrase     string
	MasterKey      *bip32.Key
}

func NewWalletKey(size SeedPhraseSize) *WalletKey {
	seedPhrase := GenerateSeedPhrase(size)
	masterKey := GenerateMasterKey(seedPhrase)

	return &WalletKey{
		SeedPhraseSize: size,
		SeedPhrase:     seedPhrase,
		MasterKey:      masterKey,
	}
}

func (w *WalletKey) GenerateSecp256k1Keys(coinType uint32) (*hdkeychain.ExtendedKey, error) {
	masterKey := w.MasterKey
	accountKey, _ := masterKey.NewChildKey(44 + hdkeychain.HardenedKeyStart)
	accountKey, _ = accountKey.NewChildKey(coinType + hdkeychain.HardenedKeyStart)
	accountKey, _ = accountKey.NewChildKey(0)
	accountKey, _ = accountKey.NewChildKey(0)

	cointMasterKey, err := hdkeychain.NewKeyFromString(accountKey.String())

	if err != nil {
		log.Fatal(err)
	}

	return cointMasterKey.Derive(0)
}

func (w *WalletKey) GenerateBitcoinKeys() BitcoinChildKeys {
	childKey, err := w.GenerateSecp256k1Keys(0)

	if err != nil {
		log.Fatal(err)
	}

	btcPublicKey, err := childKey.ECPubKey()

	compressedPubKey, err := btcutil.NewAddressPubKey(btcPublicKey.SerializeCompressed(), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	segwit, err := btcutil.NewAddressWitnessPubKeyHash(btcutil.Hash160(btcPublicKey.SerializeCompressed()), &chaincfg.MainNetParams)
	if err != nil {
		log.Fatal(err)
	}

	btcPrivateKey, err := childKey.ECPrivKey()

	if err != nil {
		log.Fatal(err)
	}

	return BitcoinChildKeys{
		PrivateKey:    btcPrivateKey.Key.String(),
		XCoordinate:   btcPublicKey.X().String(),
		YCoordinate:   btcPublicKey.Y().String(),
		P2pkhAddress:  compressedPubKey.EncodeAddress(),
		SegwitAddress: segwit.EncodeAddress(),
	}
}

func (w *WalletKey) GenerateEthereumKeys() WalletChildKeys {
	childKey, err := w.GenerateSecp256k1Keys(60)

	if err != nil {
		log.Fatal(err)
	}

	ethPublicKey, err := childKey.ECPubKey()

	if err != nil {
		log.Fatal(err)
	}

	address := crypto.PubkeyToAddress(*ethPublicKey.ToECDSA())

	return WalletChildKeys{
		PrivateKey:    childKey.String(),
		XCoordinate:   ethPublicKey.X().String(),
		YCoordinate:   ethPublicKey.Y().String(),
		PublicAddress: address.Hex(),
	}
}

func (w *WalletKey) GenerateSolanaKeys() WalletChildKeys {

	account, _ := types.AccountFromSeed([]byte(w.SeedPhrase)[:32])

	privateKey := account.PrivateKey

	privateKeyHex := hex.EncodeToString(privateKey)

	publicKey := account.PublicKey

	return WalletChildKeys{
		PrivateKey:    privateKeyHex,
		PublicAddress: publicKey.ToBase58(),
	}
}

func GenerateSeedPhrase(size SeedPhraseSize) string {
	envSeedPhrase := os.Getenv("SEED_PHRASE")
	if envSeedPhrase != "" {
		return envSeedPhrase
	}
	var entropySize int
	switch size {
	case SeedPhraseSize12:
		entropySize = 128
	case SeedPhraseSize24:
		entropySize = 256
	default:
		entropySize = 256
	}

	entropy, _ := bip39.NewEntropy(entropySize)

	seedPhrase, _ := bip39.NewMnemonic(entropy)

	return seedPhrase
}

func GenerateMasterKey(seedPhrase string) *bip32.Key {
	seed := bip39.NewSeed(seedPhrase, "")

	masterKey, _ := bip32.NewMasterKey(seed)

	return masterKey
}
