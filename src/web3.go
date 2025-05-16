package src

import (
	"context"
	"crypto/ecdsa"
	"encoding/hex"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
)

type StructMessage struct {
	Domain    string `json:"domain"`
	Address   string `json:"address"`
	ChainId   int64  `json:"chainId"`
	IssuedAt  string `json:"issuedAt"`
	Nonce     string `json:"nonce"`
	Statement string `json:"statement"`
	URI       string `json:"uri"`
	Version   string `json:"version"`
}

func GetPrivateFromHex(privateKeyHex string) *ecdsa.PrivateKey {
	if len(privateKeyHex) >= 2 && privateKeyHex[:2] == "0x" {
		privateKeyHex = privateKeyHex[2:]
	}

	privateKey, err := crypto.HexToECDSA(privateKeyHex)
	if err != nil {
		log.Fatalf("Error decoding private key: %v", err)
	}

	return privateKey
}

func PrivateToAddress(privateKeyHex string) common.Address {
	privateKey := GetPrivateFromHex(privateKeyHex)

	publicKey := privateKey.Public()
	publicKeyECDSA, ok := publicKey.(*ecdsa.PublicKey)
	if !ok {
		log.Fatal("Error type casting for public key")
	}

	address := crypto.PubkeyToAddress(*publicKeyECDSA)

	return address
}

func SignMessage(privateKeyHex string, message StructMessage) (string, error) {
	privateKeyHex = strings.TrimPrefix(privateKeyHex, "0x")

	privateKeyBytes, err := hex.DecodeString(privateKeyHex)
	if err != nil {
		return "", fmt.Errorf("failed to decode private key: %v", err)
	}

	privateKey, err := crypto.ToECDSA(privateKeyBytes)
	if err != nil {
		return "", fmt.Errorf("failed to create ECDSA private key: %v", err)
	}

	messageText := fmt.Sprintf(
		"%s wants you to sign in with your Ethereum account:\n%s\n\n%s\n\nURI: %s\nVersion: %s\nChain ID: %d\nNonce: %s\nIssued At: %s",
		message.Domain,
		message.Address,
		message.Statement,
		message.URI,
		message.Version,
		message.ChainId,
		message.Nonce,
		message.IssuedAt,
	)

	prefix := fmt.Sprintf("\x19Ethereum Signed Message:\n%d", len(messageText))
	data := []byte(prefix + messageText)
	hash := crypto.Keccak256Hash(data)

	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to sign message: %v", err)
	}

	if signature[64] < 27 {
		signature[64] += 27
	}

	return "0x" + hex.EncodeToString(signature), nil
}

func SendLegacyTransaction(client *ethclient.Client, toAddress common.Address, value *big.Int, gasLimit *uint64, callData []byte, gasPrice *big.Int, privateKeyHex string) string {
	privateKey := GetPrivateFromHex(privateKeyHex)
	fromAddress := PrivateToAddress(privateKeyHex)

	chainID, err := client.ChainID(context.Background())
	if err != nil {
		log.Fatalf("Error chainID: %v", err)
	}

	nonce, err := client.PendingNonceAt(context.Background(), fromAddress)
	if err != nil {
		log.Fatalf("Error nonce: %v", err)
	}

	if gasPrice == nil {
		gasPrice, err = client.SuggestGasPrice(context.Background())
		if err != nil {
			log.Fatalf("Error gasPrice: %v", err)
		}
	}

	txData := &types.LegacyTx{
		Nonce:    nonce,
		To:       &toAddress,
		Value:    value,
		Gas:      *gasLimit,
		GasPrice: gasPrice,
		Data:     callData,
	}

	tx := types.NewTx(txData)

	signer := types.NewEIP155Signer(chainID)
	signedTx, err := types.SignTx(tx, signer, privateKey)
	if err != nil {
		log.Fatalf("Error sign tx: %v", err)
	}

	err = client.SendTransaction(context.Background(), signedTx)
	if err != nil {
		log.Fatalf("Error send tx: %v", err)
	}

	return signedTx.Hash().String()
}
