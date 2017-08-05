package config

import (
	"github.com/spf13/viper"
	"fmt"
)

type Config struct {
	HeartbeatInterval int
	Ip string
	Port int
	RemotePeers [] string
}

func LoadConfig() *Config {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath("../")
	viper.AddConfigPath(".")

	viper.SetDefault("heartbeatinterval", 100)
	viper.SetDefault("ip", "127.0.0.1")
	viper.SetDefault("port", 9010)
	viper.SetDefault("remotepeers", []string{})

	err := viper.ReadInConfig() // Find and read the config file
	if err != nil { // Handle errors reading the config file
		panic(fmt.Errorf("Fatal error config file: %s \n", err))
	}

	return &Config{
		viper.Get("heartbeatinterval").(int),
		viper.GetString("ip"),
		viper.Get("port").(int),
		viper.GetStringSlice("remotepeers"),
	}
}
