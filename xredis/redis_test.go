package xredis

import (
	"fmt"
	"github.com/duchiporexia/goutils/xconfig"
	"github.com/go-redis/redis/v8"
	"testing"
	"time"
)

func init() {
	var cfg RedisConfig
	xconfig.LoadConfig(&cfg)
	Init(&cfg)
}

func TestExampleClient(t *testing.T) {
	err := Set("key", []byte("value"), time.Second*60)
	if err != nil {
		panic(err)
	}

	val, err := Get("key")
	if err != nil {
		panic(err)
	}
	fmt.Println("key:[", string(val), "]")

	val, err = Get("default:key")
	if err != nil {
		fmt.Printf("default:key err:%v\n", err)
	} else {
		fmt.Println("default:key:[", string(val), "]")
	}

	val2, err := Get("key2")
	if err == redis.Nil {
		fmt.Println("key2 does not exist")
	} else if err != nil {
		fmt.Printf("key2 err:%v\n", err)
	} else {
		fmt.Println("key2", val2)
	}

	results, err := MGet([]string{"dd", "key", "key2"})
	if err != nil {
		fmt.Printf("results err:%v\n", err)
	}
	for idx, result := range results {
		fmt.Printf("idx:%d, result:%v\n", idx, result)
		if result == nil {
			fmt.Printf("==> result: %v\n", result)
		}
	}

	keys, _ := Keys("key")
	fmt.Printf("keys:%v\n", keys)
	// Output: key value
	// key2 does not exist
}
