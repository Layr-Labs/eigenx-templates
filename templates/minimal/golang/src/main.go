package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

func main() {
	godotenv.Load()

	mnemonic := os.Getenv("MNEMONIC")
	if mnemonic == "" {
		log.Fatal("MNEMONIC environment variable is not set")
	}

	wallet, err := hdwallet.NewFromMnemonic(mnemonic)
	if err != nil {
		log.Fatalf("Error creating wallet from mnemonic: %v", err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		log.Fatalf("Error deriving account: %v", err)
	}

	fmt.Printf("First wallet address: %s\n", account.Address.Hex())
}
