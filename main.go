package main

import (
	"godis/config"
	tcp2 "godis/interface/tcp"
	"godis/lib/logger"
	handler2 "godis/redis/handler"
	"godis/tcp"
	"os"
)

/*
@author: shg
@since: 2023/2/22 4:11 AM
@mail: shgang97@163.com
*/

var defaultProperties = &config.ServerProperties{
	Bind:           "0.0.0.0",
	Port:           9379,
	AppendOnly:     false,
	AppendFilename: "",
	MaxClients:     1000,
}

func main() {
	logger.Setup(&logger.Settings{
		Path:       "logs",
		Name:       "godis",
		Ext:        "log",
		TimeFormat: "2006-01-02",
	})
	configFilename := os.Getenv("CONFIG")
	if configFilename == "" {
		if fileExists("redis.conf") {
			config.SetupConfig("redis.conf")
		} else {
			config.Properties = defaultProperties
		}
	} else {
		config.SetupConfig(configFilename)
	}
	cfg := &tcp.Config{
		Address:    "localhost:8080",
		MaxConnect: 100,
		Timeout:    100,
	}
	var handler tcp2.Handler // handler 接口类型
	handler = &handler2.GodisHandler{}
	err := tcp.ListenAndServeWithSignal(cfg, handler)
	if err != nil {
		logger.Error(err)
	}
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	return err == nil && !info.IsDir()
}
