package main

import (
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"

	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip32"
	"github.com/tyler-smith/go-bip39"
)

type Wallet struct {
	Index      int    `json:"index"`
	Mnemonic   string `json:"mnemonic"`
	Address    string `json:"address"`
	PrivateKey string `json:"privateKey"`
}

func createNewETHWallet(numWallets int) []Wallet {
	var wallets []Wallet

	for i := 0; i < numWallets; i++ {
		// 1. 生成助记词
		entropy, err := bip39.NewEntropy(128) // 128 bits = 12 个助记词
		if err != nil {
			log.Fatalf("Failed to generate entropy: %s", err)
		}
		mnemonic, err := bip39.NewMnemonic(entropy)
		if err != nil {
			log.Fatalf("Failed to generate mnemonic: %s", err)
		}

		// 2. 基于助记词生成种子
		seed := bip39.NewSeed(mnemonic, "")

		// 3. 使用 BIP32 派生路径（m / 44' / 60' / 0' / 0 / 0）
		masterKey, err := bip32.NewMasterKey(seed)
		if err != nil {
			log.Fatalf("Failed to generate master key: %s", err)
		}
		purposeKey, err := masterKey.NewChildKey(bip32.FirstHardenedChild + 44) // 44'
		if err != nil {
			log.Fatalf("Failed to derive purpose key: %s", err)
		}
		coinKey, err := purposeKey.NewChildKey(bip32.FirstHardenedChild + 60) // 60' (ETH)
		if err != nil {
			log.Fatalf("Failed to derive coin key: %s", err)
		}
		accountKey, err := coinKey.NewChildKey(bip32.FirstHardenedChild + 0) // 0'
		if err != nil {
			log.Fatalf("Failed to derive account key: %s", err)
		}
		changeKey, err := accountKey.NewChildKey(0) // 0 (external chain)
		if err != nil {
			log.Fatalf("Failed to derive change key: %s", err)
		}
		addressKey, err := changeKey.NewChildKey(0) // address index = 0
		if err != nil {
			log.Fatalf("Failed to derive address key: %s", err)
		}

		// 4. 提取私钥
		privateKeyBytes := addressKey.Key
		privateKey, err := crypto.ToECDSA(privateKeyBytes)
		if err != nil {
			log.Fatalf("Failed to convert to ECDSA private key: %s", err)
		}
		privateKeyHex := hex.EncodeToString(crypto.FromECDSA(privateKey))

		// 5. 使用私钥生成地址
		publicKey := privateKey.Public()
		publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
		if !ok {
			log.Fatalf("Failed to assert type: %T", publicKey)
		}
		address := crypto.PubkeyToAddress(*publicKeyECDSA).Hex()

		// 6. 构建钱包
		wallet := Wallet{
			Index:      i,
			Mnemonic:   mnemonic,
			Address:    address,
			PrivateKey: privateKeyHex,
		}
		wallets = append(wallets, wallet)
	}

	return wallets
}

func main() {
	numWallets := 10
	wallets := createNewETHWallet(numWallets)

	// 打印生成的钱包
	for _, wallet := range wallets {
		fmt.Printf("Wallet[%d]:\n", wallet.Index)
		fmt.Printf("  Mnemonic: %s\n", wallet.Mnemonic)
		fmt.Printf("  Address: %s\n", wallet.Address)
		fmt.Printf("  PrivateKey: %s\n", wallet.PrivateKey)
	}
}
