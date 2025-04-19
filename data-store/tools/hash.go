package tools

import (
	"crypto"
	"fmt"
	"log"
)

func CalculateHash(data any) ([]byte, error) {
	cryptoHash := crypto.SHA256.New()
	_, err := cryptoHash.Write([]byte(fmt.Sprintf("%v", data)))
	if err != nil {
		log.Printf("failed to calculate hash: %v\n", err)
		return nil, nil
	}
	hash := cryptoHash.Sum(nil)
	return hash, nil
}
