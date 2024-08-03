package tool

import (
	"crypto/rand"
	"math/big"
	"time"
)

type Tools struct{}

func NewTools() *Tools {
	return &Tools{}
}

func (t *Tools) GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func (t *Tools) GetTime() string {
	return time.Now().Format("2006-01-02 15:04:05")
}
