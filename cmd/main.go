package main

import (
	"encoding/json"
	"enso/src"
	ai "enso/src/enso"
	"fmt"
	"log"
	"math/big"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
)

func main() {
	data, err := src.ParseFile("data.txt")
	if err != nil {
		fmt.Printf("Error parse %v", err)
	}

	proxy, err := src.ParseFile("proxy.txt")
	if err != nil {
		fmt.Printf("Error parse %v", err)
		return
	}

	if len(proxy) < len(data) {
		fmt.Print("Not enough proxy")
		return
	}

	words, err := src.ParseFile("words.txt")
	if err != nil {
		fmt.Printf("Error parse %v", err)
	}

	for i := 0; i < len(data); i++ {
		parts := strings.Split(data[i], ":")
		if len(parts) != 2 {
			fmt.Printf("Error format 'адрес:uuid'. String: %v", i+1)
			return
		}

		privateKey := parts[0]
		address := src.PrivateToAddress(privateKey).String()
		uuid := parts[1]
		projectSlug := src.GenerateProjectSlug(words, 1)

		log.Printf("Start process address: %v", address)

		//SEARCH AI
		for n := 0; n < 5; {
			search := ai.Search("What is Uniswap v4?", "b4393b93-e603-426d-8b9f-0af145498c92", address, proxy[i])

			type Response struct {
				Answer string `json:"answer"`
			}

			var result Response
			err := json.Unmarshal([]byte(search), &result)
			if err != nil {
				fmt.Println("Error encode JSON:", err)
				return
			}

			if len(result.Answer) > 5 {
				log.Printf("%v %v AI Search", address, n+1)
				n++
			}

			src.Timeout(500)
		}

		//TRANSACTION AI
		for n := 0; n < 5; {
			rpcURL := "https://base.llamarpc.com"

			client, err := ethclient.Dial(rpcURL)
			if err != nil {
				log.Fatalf("Error connect to client: %v", err)
			}
			var gasLimit uint64 = 30000
			hashTx := src.SendLegacyTransaction(client, src.PrivateToAddress(privateKey), big.NewInt(1e14), &gasLimit, common.FromHex("0x"), big.NewInt(2e7), privateKey)
			log.Printf("Send TX: https://basescan.org/tx/%v", hashTx)
			src.Timeout(2500)

			message := src.StructMessage{
				Domain:    "enso.brianknows.org",
				Address:   address,
				ChainId:   8453,
				IssuedAt:  src.GetCurrentISOTime(),
				Nonce:     ai.GetNonce(proxy[i]),
				Statement: "By signing this message, you confirm you have read and accepted the following Terms and Conditions: https://terms.enso.build/",
				URI:       "https://enso.brianknows.org",
				Version:   "1",
			}

			sign, err := src.SignMessage(privateKey, message)
			if err != nil {
				fmt.Print("Error signature")
				return
			}

			verify, err := ai.Verify(message, sign, proxy[i])
			if err != nil {
				fmt.Print("Error verify")
				return
			}
			points := ai.Points(hashTx, "transfer", 8453, address, proxy[i], verify.Token)

			type ResponsePoints struct {
				Result string `json:"result"`
			}

			var result ResponsePoints
			err = json.Unmarshal([]byte(points), &result)
			if err != nil {
				fmt.Println("Error encode JSON:", err)
				return
			}

			if result.Result == "" {
				log.Printf("%v %v AI Transaction", address, n+1)
				n++
			}

			src.Timeout(500)
		}

		//CREATE DEX
		for n := 0; n < 5; {
			result := ai.TrackProject(projectSlug, "shortcuts-widget", address, uuid, proxy[i])
			if err != nil {
				fmt.Printf("Error TrackProjectWithProxy %s: %v\n", proxy, err)
				continue
			}

			if result == "You've earned points for creating a project!" {
				log.Printf("%v: %v. Create %v %v", address, result, n+1, projectSlug)
				n++
			}
		}
	}
}
