package controllers

import (
	"context"
	"errors"
	"os"

	"go-practice/protos"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type profile struct {
	client protos.TestRpcServiceClient
}

var Client profile
var ctx = context.Background()

func init() {
	host := getAuthHost()
	option := grpc.WithTransportCredentials(insecure.NewCredentials())

	conn, err := grpc.Dial(host, option)

	if err != nil {
		panic(err)
	}

	Client.client = protos.NewTestRpcServiceClient(conn)
}

/**
 * 写入 users 表信息
 * 把 数据 写入到 users 数据库
 */
func (p *profile) AddUserData(param string) (*protos.AddUserDataResponse, error) {
	result := &protos.AddUserDataResponse{}

	if p == nil {
		return result, errors.New("profile is nil")
	}

	req := protos.AddUserDataRequest{
		Param: param,
	}

	res, err := Client.client.AddUserData(ctx, &req)

	if err != nil {
		return result, err
	}

	return res, nil
}

/**
 * 获取 users 表信息
 * 根据 uid 与 指定的 字段集（必传）返回信息
 * 返回字符串(map[string]string)，根据需要处理 string to int...
 */
func (p *profile) GetUserData(uid, param string) (*protos.GetUserDataResponse, error) {
	result := &protos.GetUserDataResponse{}

	if p == nil {
		return result, errors.New("profile is nil")
	}

	req := protos.GetUserDataRequest{
		Uid:   uid,
		Param: param,
	}

	res, err := Client.client.GetUserData(ctx, &req)

	if err != nil {
		return result, err
	}

	return res, nil
}

func getAuthHost() string {
	res := os.Getenv("GO_MODE")

	result := "127.0.0.1:9527"

	if res != "release" {
		result = "127.0.0.1:9527"
	}

	return result
}
