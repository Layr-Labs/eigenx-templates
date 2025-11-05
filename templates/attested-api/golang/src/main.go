package main

import (
	"crypto/ecdsa"
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	hdwallet "github.com/miguelmota/go-ethereum-hdwallet"
)

type signer struct {
	address    common.Address
	privateKey *ecdsa.PrivateKey
}

type beaconResponse struct {
	RandomNumber        string `json:"randomNumber"`
	RandomNumberDecimal string `json:"randomNumberDecimal"`
	Timestamp           string `json:"timestamp"`
	Message             string `json:"message"`
	MessageHash         string `json:"messageHash"`
	Signature           string `json:"signature"`
	Signer              string `json:"signer"`
}

// newSigner derives the application's signing account from the provided mnemonic
func newSigner(mnemonic string) (*signer, error) {
	trimmed := strings.TrimSpace(mnemonic)
	if trimmed == "" {
		return nil, fmt.Errorf("mnemonic is empty")
	}

	wallet, err := hdwallet.NewFromMnemonic(trimmed)
	if err != nil {
		return nil, fmt.Errorf("wallet init: %w", err)
	}

	path := hdwallet.MustParseDerivationPath("m/44'/60'/0'/0/0")
	account, err := wallet.Derive(path, false)
	if err != nil {
		return nil, fmt.Errorf("wallet derive: %w", err)
	}

	key, err := wallet.PrivateKey(account)
	if err != nil {
		return nil, fmt.Errorf("wallet private key: %w", err)
	}

	return &signer{address: account.Address, privateKey: key}, nil
}

// buildResponse generates a random number and attests to it using the application's wallet
func buildResponse(s *signer) (*beaconResponse, error) {
	// Generate cryptographically secure random number
	buf := make([]byte, 32)
	if _, err := rand.Read(buf); err != nil {
		return nil, fmt.Errorf("entropy: %w", err)
	}

	randomNumber := "0x" + hex.EncodeToString(buf)
	randomNumberDecimal := new(big.Int).SetBytes(buf).String()
	timestamp := time.Now().UTC().Truncate(time.Millisecond).Format("2006-01-02T15:04:05.000Z")
	message := fmt.Sprintf("RandomnessBeacon|%s|%s", randomNumber, timestamp)
	hash := accounts.TextHash([]byte(message))

	// Sign the message using the application's wallet to attest to the random value
	signature, err := crypto.Sign(hash[:], s.privateKey)
	if err != nil {
		return nil, fmt.Errorf("sign: %w", err)
	}

	return &beaconResponse{
		RandomNumber:        randomNumber,
		RandomNumberDecimal: randomNumberDecimal,
		Timestamp:           timestamp,
		Message:             message,
		MessageHash:         hexutil.Encode(hash[:]),
		Signature:           hexutil.Encode(signature),
		Signer:              s.address.Hex(),
	}, nil
}

func main() {
	_ = godotenv.Load()

	mnemonic := os.Getenv("MNEMONIC")
	if mnemonic == "" {
		log.Fatal("MNEMONIC environment variable is not set")
	}

	signer, err := newSigner(mnemonic)
	if err != nil {
		log.Fatalf("failed to initialize signer: %v", err)
	}

	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())
	router.GET("/random", func(c *gin.Context) {
		resp, err := buildResponse(signer)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, resp)
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	if err := router.Run("0.0.0.0:" + port); err != nil {
		log.Fatalf("server failed: %v", err)
	}
}
