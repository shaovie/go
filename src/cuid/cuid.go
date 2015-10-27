package main

import (
	"flag"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"syscall"
	"time"

	"github.com/garyburd/redigo/redis"
	"db"
	"ilog"
	"prof"
)

// Launch args
var (
	sConfigPath = ""
	showVersion = false
	toProfile   = false
	toDaemon    = false

	redisConn redis.Conn
)

func usage() {
	var usageStr = `
        Server options:
        -c  FILE                Configuration file
        -p                      To profile
        -d                      To daemon mode

        Common options:
        -h                      Show this message
        -v                      Show version
        `
	println(usageStr)
	os.Exit(0)
}

func parseFlag() {
	flag.StringVar(&sConfigPath, "c", sConfigPath, "Configuration file.")
	flag.BoolVar(&showVersion, "v", showVersion, "Show version.")
	flag.BoolVar(&toProfile, "p", toProfile, "To profile.")
	flag.BoolVar(&toDaemon, "d", toDaemon, "To daemon mode.")

	flag.Usage = usage
	flag.Parse()
}

func signalHandle() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT)
	for {
		sig := <-ch
		if sig == syscall.SIGHUP {
			ilog.Rinfo("signal hup")
		} else if sig == syscall.SIGINT {
			ilog.Rinfo("signal int")
		}
	}
}

func redisTest(ch chan<- int, key string) {
	redisConn, err := redis.DialTimeout("tcp", "172.18.8.24:6001", 1*time.Second, 1*time.Second, 1*time.Second)
	if err != nil {
		ilog.Error("redis - " + err.Error())
		os.Exit(0)
	}
	redisConn.Do("DEL", "hello"+key)
	for i := 0; i < 1000000; i++ {
		_, err := redisConn.Do("INCR", "hello"+key)
		if err != nil {
			ilog.Error("redis - " + err.Error())
		}
	}
	ret, err := redisConn.Do("GET", "hello"+key)
	ilog.Rinfo("redis - " + fmt.Sprintf("%s", ret))
	ch <- 1
}

func doRedisTest() {
	closeCh := make(chan int)
	for i := 0; i < 3; i++ {
		go redisTest(closeCh, strconv.Itoa(i))
	}
	for i := 0; i < 3; i++ {
		<-closeCh
	}
	os.Exit(0)
}

func main() {
	parseFlag()

	if showVersion {
		println("version 1.0.0")
		os.Exit(0)
	}

	if err := loadConfig(); err != nil {
		println("config -", err.Error())
		os.Exit(0)
	}

	if toDaemon {
		if err := daemon(); err != nil {
			println("daemon -", err.Error())
			os.Exit(0)
		}
	}

	go signalHandle()

	if toProfile {
		prof.StartProf()
	}

	if err := initSvc(); err != nil {
		println("init -", err.Error())
		os.Exit(1)
	}

	// output pid
	if err := outputPid(ServerConfig.PidFile); err != nil {
		ilog.Error("pid - " + err.Error())
		os.Exit(0)
	}

	ilog.Rinfo("launch ok")
	runtime.GOMAXPROCS(runtime.NumCPU())

	//doRedisTest()
	doMysqlTest()

	index := func(w http.ResponseWriter, r *http.Request) {
		ret, err := redisConn.Do("PING")
		if err != nil {
			ilog.Error("redis - " + err.Error())
		} else {
			io.WriteString(w, fmt.Sprintf("%s", ret))
		}
	}
	http.HandleFunc("/", index)
	http.ListenAndServe(":8081", nil)

	/*
		listener, err := net.Listen("tcp", "0.0.0.0:8088")
		if err != nil {
			ilog.Error("listen fail " + err.Error())
			os.Exit(0)
		}
		defer listener.Close()

		for {
			conn, err := listener.Accept()
			if err != nil {
				ilog.Error("accept error " + err.Error())
				break
			}
			go func(conn net.Conn) {
				for {
					buf := make([]byte, 1024)
					_, err := conn.Read(buf)
					if err != nil {
						conn.Close()
						return
					}

					_, err = conn.Write([]byte("HTTP/1.1 200 OK\r\nConnection: keep-alive\r\nContent-Length: 5\r\n\r\nhello"))
				}
			}(conn)
		}
	*/

	os.Exit(1)
}
