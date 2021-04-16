package tool

import (
	"crypto/rand"
	"encoding/json"
	"go-practice/libs/types"
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

func UnmarshalJson(date string) types.MapStringInterface {
	var res types.MapStringInterface

	_ = json.Unmarshal([]byte(date), &res)

	return res
}

func GetRandmod(length int) int {
	res, err := rand.Int(rand.Reader, big.NewInt(int64(length)))

	if err != nil {
		log.Fatal(err)
	}

	return int(res.Int64())
}

func GetGinMODE() string {
	res := os.Getenv("GIN_MODE")

	if res != "release" {
		res = "test"
	}

	return res
}
