package main

import (
	"fmt"
	wallets "go-blockchain/wallets"
	"path/filepath"

	"github.com/joho/godotenv"
)

func main() {
	path, _ := filepath.Abs("../.env")
	err := godotenv.Load(path)

	if err != nil {
		panic("Error loading .env file")
	}
	seedPhrase := wallets.GenerateSeedPhrase(wallets.SeedPhraseSize24)
	generateMasterKey := wallets.GenerateMasterKey(seedPhrase)

	fmt.Println(seedPhrase)
	fmt.Println(generateMasterKey)

	walletKey := wallets.NewWalletKey(wallets.SeedPhraseSize24)
	bitcoinKey := walletKey.GenerateBitcoinKeys()

	fmt.Println("bitcoinKey %+v\n", bitcoinKey)

	// 	ethereumKey := walletKey.GenerateEthereumKeys()

	// 	fmt.Println("ethereumKey %+v\n", ethereumKey)

	// walletKey.GenerateSolanaKeys()
}
