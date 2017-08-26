package config

import (
	"github.com/spf13/viper"
	"fmt"
	"path/filepath"
)

type Config struct {
	HeartbeatInterval int
	Address string
	RemotePeers [] string
}

func LoadConfig(config string) *Config {
	path := filepath.Dir(config)
	filename := filepath.Base(config)
	conf := filename[0:len(filename) - len(filepath.Ext(filename))]

	viper.SetConfigName(conf)
	viper.SetConfigType("yaml")
	viper.AddConfigPath(path)

	viper.SetDefault("heartbeatinterval", 100)
	viper.SetDefault("ip", "127.0.0.1:9010")
	viper.SetDefault("remotepeers", []string{})

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return &Config {
		viper.Get("heartbeatinterval").(int),
		viper.GetString("address"),
		viper.GetStringSlice("remotepeers"),
	}
}
