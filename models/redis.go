package models

import (
	"context"

	"go-practice/libs/tool"

	log "github.com/sirupsen/logrus"
)

var ctx = context.TODO()

func GetUserName() (result string, err error) {
	pool := tool.GetRedisClient("users.master")
	back, err := pool.HGet(ctx, "u:644", "name").Result()

	if err != nil {
		log.Error(err)
	}

	result = back

	return result, err
}
