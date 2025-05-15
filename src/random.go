package src

import (
	"fmt"
	"math/big"
	"math/rand"
	"strings"
)

func GenerateProjectSlug(wordsList []string, numWords int) string {
	if len(wordsList) == 0 || numWords <= 0 {
		return ""
	}

	var selectedWords []string
	for i := 0; i < numWords && i < len(wordsList); i++ {
		randomIndex := rand.Intn(len(wordsList))
		selectedWords = append(selectedWords, wordsList[randomIndex])

		wordsList = append(wordsList[:randomIndex], wordsList[randomIndex+1:]...)

		if len(wordsList) == 0 {
			break
		}
	}

	randomNum, err := GenerateRandomAmount(big.NewInt(100), big.NewInt(99999), 0)
	if err != nil {
		fmt.Printf("Error generate random number: %v", err)
		return ""
	}

	wordsStr := strings.Join(selectedWords, "-")
	slug := fmt.Sprintf("%s%d.widget", wordsStr, randomNum)

	return slug
}
