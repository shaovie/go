package main

type serverConfig struct {
	LogFile string
	PidFile string
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