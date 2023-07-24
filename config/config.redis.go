package config

import "go-practice/libs/types"

var redisConfig = map[string]types.ConfRedis{
	"default-test": {
		Master:   []string{"127.0.0.1:6379"},
		Password: "",
		Db:       0,
	},
	"default-release": {
		Master:   []string{"127.0.0.1:6379"},
		Password: "",
		Db:       0,
	},
}

func GetRedisConfig() types.FullConfRedis {
	mode := getMode()
	result := types.FullConfRedis{}

	data := []string{
		"default",
	}

	for _, v := range data {
		result[v] = redisConfig[v+"-"+mode]
	}

	return result
}
