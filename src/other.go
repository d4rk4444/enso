package src

import (
	"bufio"
	"fmt"
	"math/big"
	"math/rand"
	"os"
	"strings"
	"time"
)

func GenerateRandomAmount(min *big.Int, max *big.Int, decimal int) (*big.Int, error) {
	if min.Cmp(max) >= 0 {
		return nil, fmt.Errorf("min < max")
	}

	minFloat := new(big.Float).SetInt(min)
	maxFloat := new(big.Float).SetInt(max)

	diff := new(big.Float).Sub(maxFloat, minFloat)

	randFloat := new(big.Float).SetFloat64(rand.Float64())

	randDiff := new(big.Float).Mul(randFloat, diff)

	result := new(big.Float).Add(minFloat, randDiff)

	multiplier := new(big.Int).Exp(
		big.NewInt(10),
		big.NewInt(int64(decimal)),
		nil,
	)
	multiplierFloat := new(big.Float).SetInt(multiplier)

	resultScaled := new(big.Float).Mul(result, multiplierFloat)

	resultInt := new(big.Int)
	resultScaled.Int(resultInt)

	return resultInt, nil
}

func Timeout(ms int) {
	time.Sleep(time.Duration(ms) * time.Millisecond)
}

func ParseFile(file string) ([]string, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, fmt.Errorf("Error open file: %v", err)
	}
	defer f.Close()

	scanner := bufio.NewScanner(f)
	var lines []string

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.ReplaceAll(line, "\r", "")
		line = strings.ReplaceAll(line, " ", "")

		if line != "" {
			lines = append(lines, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Error read file: %v", err)
	}

	return lines, nil
}

func GetCurrentISOTime() string {
	currentTime := time.Now().UTC()
	isoTime := currentTime.Format("2006-01-02T15:04:05.000Z")

	return isoTime
}
