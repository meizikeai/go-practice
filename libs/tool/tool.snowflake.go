package tool

import (
	"github.com/bwmarrin/snowflake"
)

type Node struct {
	cli *snowflake.Node
}

var Snowflake Node

func init() {
	node, err := snowflake.NewNode(1)

	if err != nil {
		panic(err)
	}

	Snowflake.cli = node
}

func (a *Node) GenerateSnowflake() string {
	result := Snowflake.cli.Generate()
	return result.String()
}
