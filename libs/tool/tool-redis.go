package tool

import (
	"context"
	"encoding/json"
	"go-practice/libs/types"
	"io/ioutil"
	"os"
	"strings"

	"github.com/go-redis/redis/v8"
	log "github.com/sirupsen/logrus"
)

var ctx = context.Background()
var fullDbRedis map[string]*redis.Client
var redisConfig types.FullConfRedis

func HandleLocalRedisConfig() {
	var config types.FullConfRedis

	pwd, _ := os.Getwd()
	mode := GetGinMODE()

	address := strings.Join([]string{pwd, "/conf/", mode, "-redis.json"}, "")

	res, err := ioutil.ReadFile(address)

	if err != nil {
		log.Fatal(err)
	}

	err = json.Unmarshal(res, &config)

	if err != nil {
		log.Fatal(err)
	}

	redisConfig = config
}

func GetRedisClient(key string) *redis.Client {
	return fullDbRedis[key]
}

func HandleRedisClient() {
	config := make(map[string]*redis.Client)

	one := getZookeeperRedisConfig()
	two := getLocalRedisConfig()

	for key, val := range one {
		m := getRedisConfig(val)

		config[key+".master"] = createRedisClient(m)
	}

	for key, val := range two {
		m := getRedisConfig(val)

		config[key+".master"] = createRedisClient(m)
	}

	fullDbRedis = config
}

func getLocalRedisConfig() types.FullConfRedis {
	return redisConfig
}

func createRedisClient(config types.OutConfRedis) *redis.Client {
	// db *redis.Client
	// 可长期存在，不建议频繁开启/关闭
	// defer db.Close()

	// create client
	db := redis.NewClient(&redis.Options{
		Addr:     config.Addr,
		Password: config.Password,
		DB:       config.Db,
	})

	// test client
	_, err := db.Ping(ctx).Result()

	if err != nil {
		log.Fatal(err)
	}

	return db
}

func getRedisConfig(config types.ConfRedis) types.OutConfRedis {
	where := config.Master

	i := GetRandmod(len(where))

	r := types.OutConfRedis{
		Addr:     where[i],
		Password: config.Password,
		Db:       config.Db,
	}

	return r
}
