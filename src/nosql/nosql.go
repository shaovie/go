package nosql

import (
	"errors"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var keyPrefix string
var redisConn redis.Conn
var nosqlMtx sync.Mutex
var nosqlAddr string

const maxBadConnRetries = 2

func Init(host, port, prefix string) (err error) {
	nosqlMtx.Lock()
	defer nosqlMtx.Unlock()

	keyPrefix = prefix
	nosqlAddr = host + ":" + port

	redisConn, err = open(nosqlAddr)
	return err
}

func open(addr string) (redis.Conn, error) {
	return redis.DialTimeout("tcp",
		addr,
		1*time.Second,
		1*time.Second,
		1*time.Second)
}

func doCmd(cmd string, args ...interface{}) (ret interface{}, err error) {
	nosqlMtx.Lock()
	defer nosqlMtx.Unlock()

	if redisConn == nil {
		return nil, errors.New("nosql: redis do not init")
	}

	for i := 0; i < maxBadConnRetries; i++ {
		ret, err = redisConn.Do(cmd, args)
		if err == nil {
			break
		}
		redisConn, err = open(nosqlAddr)
	}
	return ret, err
}

func Get(key string) (ret interface{}, err error) {
	ret, err = doCmd("GET", keyPrefix+key)
	return ret, err
}

func Set(key string, val interface{}) (err error) {
	_, err = doCmd("SET", keyPrefix+key, val)
	return err
}

func Del(key string) (err error) {
	_, err = doCmd("DEL", keyPrefix+key)
	return err
}

func Incr(key string) (err error) {
	_, err = doCmd("INCR", keyPrefix+key)
	return err
}
