package main

import (
	"math/rand"
	"time"

	"ilog"
)

func initSvc() (err error) {
	if err = ilog.Open(ServerConfig.LogFile); err != nil {
		return err
	}

	// Init randdom generater
	rand.Seed(time.Now().UnixNano())
	return err
}
