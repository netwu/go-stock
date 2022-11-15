package redisUtil

import (
	"fmt"
	"time"

	"github.com/go-redis/redis"
	"github.com/goravel/framework/facades"
)

type RedisUtil struct {
	client *redis.Client
}

func NewRedisUtil() *RedisUtil {
	return &RedisUtil{
		InitRedisClient(),
	}
}
func InitRedisClient() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", facades.Config.GetString("database.redis.default.host"), facades.Config.GetString("database.redis.default.port")),
		Password: facades.Config.GetString("database.redis.default.password"),
		DB:       facades.Config.GetInt("database.redis.default.database"),
	})
	_, err := client.Ping().Result()
	if err != nil {
		fmt.Println("Redis Client Error ")
	} else {
		fmt.Println("Redis Client Success")

	}
	return client

}

func (RedisUtil *RedisUtil) BatchPushRedis(key string, values []string) {
	for _, v := range values {
		err := RedisUtil.client.LPush(key, v).Err()
		if err != nil {
			panic(err)
		}
	}
}
func (RedisUtil *RedisUtil) PushRedis(key string, v string) bool {

	err := RedisUtil.client.LPush(key, v).Err()
	if err != nil {
		panic(err)
	}
	return true
}
func (RedisUtil *RedisUtil) DelKeyRedis(key string) bool {

	err := RedisUtil.client.Del(key).Err()
	if err != nil {
		panic(err)
	}
	return true
}
func (RedisUtil *RedisUtil) RedisRPop(key string) string {
	val, err := RedisUtil.client.RPop(key).Result()
	if err != nil {
		return ""
	}
	return val
}
func (RedisUtil *RedisUtil) RedisLLen(key string) int64 {
	val, err := RedisUtil.client.LLen(key).Result()
	if err != nil {
		panic(err)
	}

	return val
}
func (RedisUtil *RedisUtil) RedisSet(key string, val string, expire int) {

	val, err := RedisUtil.client.Set(key, val, time.Duration(expire)*time.Second).Result()
	if err != nil {
		panic(err)
	}

}
func (RedisUtil *RedisUtil) RedisGet(key string) string {
	val, err := RedisUtil.client.Get(key).Result()
	if err != nil {
		return ""
	}
	return val

}
func (RedisUtil *RedisUtil) RedisBatchPop(key string, count int) []string {

	var ret []string
	for i := 0; i < count; i++ {
		ret = append(ret, RedisUtil.RedisRPop(key))
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

func (RedisUtil *RedisUtil) RedisInc(key string) {
	RedisUtil.client.Incr(key)
}
