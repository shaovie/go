package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/garyburd/redigo/redis"
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
			ilog.Rinfo("signal - receive hup")
		} else if sig == syscall.SIGINT {
			ilog.Rinfo("signal - receive int")
		}
	}
}

func main() {
	parseFlag()

	if showVersion {
		println("version 1.0.0")
		os.Exit(0)
	}

	if err := loadConfig(); err != nil {
		println("config - ", err.Error())
		os.Exit(0)
	}

	if toDaemon {
		if err := daemon(); err != nil {
			println("daemon - ", err.Error())
			os.Exit(0)
		}
	}

	runtime.GOMAXPROCS(runtime.NumCPU() * 2)

	go signalHandle()

	if toProfile {
		prof.StartProf()
	}

	if err := initSvc(); err != nil {
		println("init - ", err.Error())
		os.Exit(1)
	}

	if err := outputPid(ServerConfig.PidFile); err != nil {
		ilog.Error("pid - " + err.Error())
		os.Exit(0)
	}

	ilog.Rinfo("launch ok")

	err := startServer()
	if err != nil {
		println("start - ", err.Error())
	}

	os.Exit(1)
}
