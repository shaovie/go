package main

import (
	"math/rand"
	"time"

	"db"
	"ilog"
	"mc"
	"nosql"
)

func initSvc() (err error) {
	if err = ilog.Open(ServerConfig.LogFile); err != nil {
		return err
	}

	// Init randdom generater
	rand.Seed(time.Now().UnixNano())

	err = mc.Init(ServerConfig.MC.Host, ServerConfig.MC.Port, "mc")
	if err != nil {
		ilog.Error("mc - init fail! " + err.Error())
	}

	err = nosql.Init(ServerConfig.Nosql.Host, ServerConfig.Nosql.Port, "nosql")
	if err != nil {
		ilog.Error("nosql - init fail! " + err.Error())
	}

	err = db.Init(ServerConfig.DB.Host,
		ServerConfig.DB.Port,
		ServerConfig.DB.User,
		ServerConfig.DB.Passwd,
		ServerConfig.DB.Charset,
		ServerConfig.DB.Name)
	if err != nil {
		ilog.Error("db - init fail! " + err.Error())
	}

	return err
}
