package nosql

import (
	"fmt"
	"os"
	"strconv"
	"testing"
	"time"

	"github.com/garyburd/redigo/redis"
)

func redisTest(ch chan<- int, key string) {
	redisConn, err := redis.DialTimeout("tcp", "172.18.8.24:6379", 1*time.Second, 1*time.Second, 1*time.Second)
	if err != nil {
		fmt.Println("redis - " + err.Error())
		os.Exit(0)
	}
	redisConn.Do("DEL", "hello"+key)
	for i := 0; i < 1000000; i++ {
		_, err := redisConn.Do("INCR", "hello"+key)
		if err != nil {
			fmt.Println("redis - " + err.Error())
		}
	}
	ret, err := redisConn.Do("GET", "hello"+key)
	fmt.Println("redis - " + fmt.Sprintf("%s", ret))
	ch <- 1
}

func BenchmarkRedis(b *testing.B) {
	closeCh := make(chan int)
	for i := 0; i < 3; i++ {
		go redisTest(closeCh, strconv.Itoa(i))
	}
	for i := 0; i < 3; i++ {
		<-closeCh
	}
}
