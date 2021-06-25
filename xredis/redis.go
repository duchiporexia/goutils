package xredis

import (
	"context"
	"github.com/duchiporexia/goutils/xmsg"
	"github.com/go-redis/redis/v8"
	"time"
)

var rdb *redis.Client
var prefix = ""
var ctxBg = context.Background()

type RedisConfig struct {
	Url    string `yaml:"url" env:"URL" env-default:"localhost:6379"`
	Prefix string `yaml:"prefix" env:"PREFIX" env-default:"bkgo:"`
}

func Init(cfg *RedisConfig) {
	rdb = redis.NewClient(&redis.Options{
		Addr:     cfg.Url,
		Password: "", // no password set
		DB:       0,  // use default DB
	})
	prefix = cfg.Prefix
}

func wrapKey(key string) string {
	return prefix + key
}

func Set(key string, value []byte, expiration time.Duration) error {
	return rdb.Set(ctxBg, wrapKey(key), value, expiration).Err()
}

func SetMsg(key string, msg xmsg.MsgMarshaler, expiration time.Duration) error {
	bytes, err := msg.MarshalMsg(nil)
	if err != nil {
		return err
	}
	return Set(key, bytes, expiration)
}

func Get(key string) ([]byte, error) {
	return rdb.Get(ctxBg, wrapKey(key)).Bytes()
}

func GetMsg(key string, msg xmsg.MsgUnmarshaler) error {
	bytes, err := Get(key)
	if err != nil {
		return err
	}
	_, err = msg.UnmarshalMsg(bytes)
	return err
}

func MGetRaw(rawKeys []string) ([]interface{}, error) {
	return rdb.MGet(ctxBg, rawKeys...).Result()
}

func MGet(keys []string) ([]interface{}, error) {
	rawKeys := make([]string, len(keys))
	for i, key := range keys {
		rawKeys[i] = wrapKey(key)
	}
	return rdb.MGet(ctxBg, rawKeys...).Result()
}

func Keys(pattern string) ([]string, error) {
	return rdb.Keys(ctxBg, wrapKey(pattern)).Result()
}

func Del(keys ...string) error {
	newKeys := make([]string, len(keys))
	for i := 0; i < len(keys); i++ {
		newKeys[i] = wrapKey(keys[i])
	}
	return rdb.Del(ctxBg, newKeys...).Err()
}
