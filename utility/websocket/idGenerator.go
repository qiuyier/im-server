package websocket

import "github.com/bwmarrin/snowflake"

func IdGen() int64 {
	snowflakeGen, err := snowflake.NewNode(1)
	if err != nil {
		panic(err)
	}

	return snowflakeGen.Generate().Int64()
}
