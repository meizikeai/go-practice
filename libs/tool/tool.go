package tool

import (
	"crypto/rand"
	"encoding/json"
	"math/big"
	"os"

	_ "github.com/go-sql-driver/mysql"
	log "github.com/sirupsen/logrus"
)

func MarshalJson(date interface{}) []byte {
	res, err := json.Marshal(date)

	if err != nil {
		log.Error(err)
	}

	return res
}

func UnmarshalJson(date string) map[string]interface{} {
	var res map[string]interface{}

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func GetRandmod(length int) int64 {
	result := int64(0)
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		return result
	}

	return res.Int64()
}

func GetMODE() string {
	res := os.Getenv("GIN_MODE")

	if res != "release" {
		res = "test"
	}

	return res
}
