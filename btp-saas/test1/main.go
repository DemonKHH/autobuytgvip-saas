package main

import (
	"crypto/ed25519"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/tonkeeper/tongo/liteapi"
	"github.com/tonkeeper/tongo/wallet"
)

func GenerateWalletInfo() {
	seedStr := "speak write argue ordinary shy melody curve mansion receive cupboard round climb aspect second region wise stuff act improve slim basket ability crystal pretty"
	client, err := liteapi.NewClientWithDefaultMainnet()
	if err != nil {
		log.Printf("Unable to create lite client: %v\n", err)
		return
	}

	w, err := wallet.DefaultWalletFromSeed(seedStr, client)
	if err != nil {
		log.Printf("Unable to create wallet: %v\n", err)
		return
	}
	addressInfo := w.GetAddress()
	fmt.Printf("address: %s\n", addressInfo.ToRaw())
	workchain := addressInfo.Workchain
	fmt.Printf("workchain: %d\n", workchain)
	privateKey, err := wallet.SeedToPrivateKey(seedStr)
	if err != nil {
		log.Printf("Unable to create private key: %v\n", err)
		return
	}
	fmt.Printf("private key: %v\n", hex.EncodeToString(privateKey))
	publicKey := privateKey.Public().(ed25519.PublicKey)
	fmt.Printf("public key: %v\n", hex.EncodeToString(publicKey))
}

func main() {
	GenerateWalletInfo()
}
