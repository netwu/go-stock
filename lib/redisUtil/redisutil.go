package redisUtil

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/goravel/framework/support/facades"
)

func InitRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", facades.Config.GetString("database.redis.default.host"), facades.Config.GetString("database.redis.default.port")),
		Password: facades.Config.GetString("database.redis.default.password"),
		DB:       facades.Config.GetInt("database.redis.default.database"),
	})
	return client
}

func BatchPushRedis(client *redis.Client, key string, values []string) {
	for _, v := range values {
		err := client.LPush(key, v).Err()
		if err != nil {
			panic(err)
		}
	}
}
func PushRedis(client *redis.Client, key string, v string) bool {
	err := client.LPush(key, v).Err()
	if err != nil {
		panic(err)
	}
	return true
}
func DelKeyRedis(client *redis.Client, key string) bool {
	err := client.Del(key).Err()
	if err != nil {
		panic(err)
	}
	return true
}
func RedisRPop(client *redis.Client, key string) string {
	val, err := client.RPop(key).Result()
	if err != nil {
		return ""
	}
	return val
}
func RedisLLen(client *redis.Client, key string) int64 {
	val, err := client.LLen(key).Result()
	if err != nil {
		panic(err)
	}

	return val
}
func RedisSet(client *redis.Client, key string, val string, expire int) {
	val, err := client.Set(key, val, time.Duration(expire)*time.Second).Result()
	if err != nil {
		panic(err)
	}

}
func RedisGet(client *redis.Client, key string) string {
	val, err := client.Get(key).Result()
	if err != nil {
		return ""
	}
	return val

}
func RedisBatchPop(client *redis.Client, key string, count int) []string {
	var ret []string
	for i := 0; i < count; i++ {
		ret = append(ret, RedisRPop(client, key))
	}
	return ret
	// pipe := client.Pipeline()
	// commands := []*redis.StringStringMapCmd{}
	// for i := 0; i < count; i++ {
	// 	commands = append(commands, pipe.RPop(ctx, key).Result())
	// }

	// pipe.Exec(ctx)

	// for _, cmd := range commands {
	// 	result, _ := cmd.Result()
	// 	fmt.Println(result)
	// 	// do something with the result
	// }
	// fmt.Println(result.Val())

}
