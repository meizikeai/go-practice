package tool

// github.com/rs/xid
import "github.com/bwmarrin/snowflake"

type snow struct {
	cli *snowflake.Node
}

var Snowflake snow

func init() {
	node, err := snowflake.NewNode(1)

	if err != nil {
		panic(err)
	}

	Snowflake.cli = node
}

func (a *snow) GenerateSnowflake() string {
	result := Snowflake.cli.Generate()
	return result.String()
}
