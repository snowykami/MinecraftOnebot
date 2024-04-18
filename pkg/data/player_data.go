package data

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
)

var ctx = context.Background()
var RedisClient = redis.NewClient(&redis.Options{
	Addr:     "localhost:6379",
	Password: "", // no password set
	DB:       0,  // use default DB
})

const namespace = "MCOB"

func GetPlayerId(playerName string) int64 {
	// 先判断是否存在，不存在生成一个新的playerId
	dbKey := namespace + ":player:" + playerName
	playerId, err := RedisClient.Get(ctx, dbKey).Int64()
	if err != nil {
		playerId = 0
	}
	if playerId == 0 {
		// 生成新的playerId
		playerId = RedisClient.Incr(ctx, "player_id").Val()
		// 保存到redis
		RedisClient.Set(ctx, dbKey, playerId, 0)
	}
	return playerId
}

func StorageMessage(messageId int64, rawMessage string) {
	// 保存消息
	RedisClient.Set(ctx, namespace+":message:"+strconv.FormatInt(messageId, 10), rawMessage, 0)
}

func GetMessage(messageId int64) string {
	// 获取消息
	return RedisClient.Get(ctx, namespace+":message:"+strconv.FormatInt(messageId, 10)).Val()
}
