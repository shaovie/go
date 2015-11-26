package main

import (
	"encoding/json"
	"errors"
	"io/ioutil"
)

type serverConfig struct {
	LogFile string
	PidFile string

	MC struct {
		Host string
		Port string
	}

	Nosql struct {
		Host string
		Port string
	}

	DB struct {
		Host    string
		Port    string
		User    string
		Passwd  string
		Charset string
		Name    string
	}
}

var ServerConfig serverConfig

func loadConfig() error {
	if sConfigPath == "" {
		return errors.New("not configure config file!")
	}

	content, err := ioutil.ReadFile(sConfigPath)
	if err != nil {
		return errors.New("read config file failed! " + err.Error())
	}

	err = json.Unmarshal(content, &ServerConfig)
	if err != nil {
		return errors.New("parse config file failed! " + err.Error())
	}
	return nil
}
