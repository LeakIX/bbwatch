package main

import (
	"github.com/LeakIX/bbwatch/app"
	"github.com/LeakIX/bbwatch/conf"
	"gopkg.in/yaml.v3"
	"os"
)

func main() {
	configFile, err := os.Open("config.yaml")
	if err != nil {
		panic(err)
	}
	var config conf.BBWatchConfig
	err = yaml.NewDecoder(configFile).Decode(&config)
	if err != nil {
		panic(err)
	}
	bbwatcher := app.NewBBWatcher(&config)
	bbwatcher.Start()
}
