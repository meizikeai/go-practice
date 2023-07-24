package config

import "go-practice/libs/types"

var mysqlConfig = map[string]types.ConfMySQL{
	"default-test": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306"},
		Username: "test",
		Password: "yintai@123",
		Database: "test",
	},
	"default-release": {
		Master:   []string{"127.0.0.1:3306"},
		Slave:    []string{"127.0.0.1:3306", "127.0.0.1:3306"},
		Username: "test",
		Password: "yintai@123",
		Database: "test",
	},
}

func GetMySQLConfig() types.FullConfMySQL {
	mode := getMode()
	result := types.FullConfMySQL{}

	data := []string{
		"default",
	}

	for _, v := range data {
		result[v] = mysqlConfig[v+"-"+mode]
	}

	return result
}
