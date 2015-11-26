package mc

import (
	"errors"
	"sync"
	"time"

	"github.com/garyburd/redigo/redis"
)

var keyPrefix string
var redisConn redis.Conn
var mcMtx sync.Mutex
var mcAddr string

const maxBadConnRetries = 2

func Init(host, port, prefix string) (err error) {
	mcMtx.Lock()
	defer mcMtx.Unlock()

	keyPrefix = prefix
	mcAddr = host + ":" + port

	redisConn, err = open(mcAddr)
	return err
}

func open(addr string) (redis.Conn, error) {
	return redis.DialTimeout("tcp",
		addr,
		1*time.Second,
		1*time.Second,
		1*time.Second)
}

func getMc() redis.Conn {
	mcMtx.Lock()
	defer mcMtx.Unlock()

	if redisConn == nil {
		return nil
	}
	if redisConn.Err() != nil {
		conn, err := open(mcAddr)
		if err == nil {
			redisConn = conn
		}
	}
	return redisConn
}
func doCmd(cmd string, args ...interface{}) (ret interface{}, err error) {
	mcMtx.Lock()
	defer mcMtx.Unlock()

	if redisConn == nil {
		return nil, errors.New("mc: redis do not init")
	}

	for i := 0; i < maxBadConnRetries; i++ {
		ret, err = redisConn.Do(cmd, args)
		if err == nil {
			break
		}
		redisConn, err = open(mcAddr)
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

/*
import (
	"time"
	"errors"
    "sync"

	"github.com/garyburd/redigo/redis"
)

var keyPrefix string
var mcMtx sync.Mutex
const maxBadConnRetries = 2

type mcConn struct {
	host string
	port string

	conn redis.Conn
}

var mc *mcConn

func getMc(key string) *mcConn {
	return mc
}

func Init(host, port, prefix string) (err error) {
    mcMtx.Lock()
    defer mcMtx.Unlock()

	keyPrefix = prefix

	mc = &mcConn{
		host: host,
		port: port,
	}
	return mc.open()
}

func (mc *mcConn) open() error {
	conn, err := redis.DialTimeout("tcp",
		mc.host+":"+mc.port,
		1*time.Second, // connect timeout
		1*time.Second, // read timeout
		1*time.Second) // write timeout

	if err == nil {
		mc.conn = conn
	}
	return err
}

func (mc *mcConn) close() {
	mc.conn.Close()
}

func doCmd(cmd string, args ...interface{}) (ret interface{}, err error) {
	mc := getMc("x")
	if mc == nil {
		return nil, errors.New("mc: get mc fail")
	}

	for i := 0; i < maxBadConnRetries; i++ {
		if mc.conn == nil {
			err = mc.open()
			if err != nil {
				continue
			}
		}
		ret, err = mc.conn.Do(cmd, args)
		if err != nil {
			mc.open()
		} else {
            break
        }
	}
	return ret, err
}

func Get(key string) (ret interface{}, err error) {
    mcMtx.Lock()
    defer mcMtx.Unlock()

	ret, err = doCmd("GET", keyPrefix+key)
    return ret, err
}

func Set(key string, val interface{}) (err error) {
    mcMtx.Lock()
    defer mcMtx.Unlock()

	_, err = doCmd("SET", keyPrefix+key, val)
    return
}

func Del(key string) (err error) {
    mcMtx.Lock()
    defer mcMtx.Unlock()

	_, err = doCmd("DEL", keyPrefix+key)
    return
}

func Incr(key string) (err error) {
    mcMtx.Lock()
    defer mcMtx.Unlock()

	_, err = doCmd("INCR", keyPrefix+key)
    return
}
*/
