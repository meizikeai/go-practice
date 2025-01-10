package tool

import (
	"crypto/rand"
	"database/sql"
	"math/big"
	"time"

	"go-practice/config"

	"github.com/go-redis/redis/v8"
)

var (
	dbMySQLCache map[string][]*sql.DB
	dbRedisCache map[string][]*redis.Client
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

// mysql
func (t *Tools) GetMySQLClient(key string) *sql.DB {
	result := dbMySQLCache[key]
	index := t.GetRandmod(len(result))

	return result[index]
}

func (t *Tools) HandleMySQLClient() {
	config := config.GetMySQLConfig()
	result := NewMySQLClient(config)

	dbMySQLCache = result.Client

	t.Stdout("MySQL is Connected")
}

func (t *Tools) CloseMySQL() {
	for _, val := range dbMySQLCache {
		for _, v := range val {
			v.Close()
		}
	}

	t.Stdout("MySQL is Close")
}

// redis
func (t *Tools) GetRedisClient(key string) *redis.Client {
	result := dbRedisCache[key]
	index := t.GetRandmod(len(result))

	return result[index]
}

func (t *Tools) HandleRedisClient() {
	config := config.GetRedisConfig()
	result := NewRedisClient(config)

	dbRedisCache = result.Client

	t.Stdout("Redis is Connected")
}

func (t *Tools) CloseRedis() {
	for _, val := range dbRedisCache {
		for _, v := range val {
			v.Close()
		}
	}

	t.Stdout("Redis is Close")
}
