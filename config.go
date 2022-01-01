package main

import (
	"github.com/mcuadros/go-defaults"
	"github.com/spf13/viper"
	"log"
)

// Server defines the object of a Game server
type Server struct {
	Name    string `mapstructure:"name"`
	Address string `mapstructure:"address"`
	Labels  map[string]interface{}
}

type loggingConfig struct {
	FileEnabled       bool   `mapstructure:"file_enabled"`
	ConsoleEnabled    bool   `mapstructure:"console_enabled"`
	FileLogsDirectory string `mapstructure:"file_logs_dir"`
}

type config struct {
	Servers        []Server `mapstructure:"servers"`
	Logging        loggingConfig
	MaxPacketNum   int   `mapstructure:"max_packet_num" default:"20"`
	MinPacketNum   int   `mapstructure:"min_packet_num" default:"4"`
	PingTimeout    int64 `mapstructure:"ping_timeout" default:"30"`  // in seconds
	PingInterval   int64 `mapstructure:"ping_interval" default:"30"` // in seconds
	WorkerPoolSize int   `mapstructure:"worker_pool_size" default:"5"`
	UIEnabled      bool  `mapstructure:"ui_enabled" default:"true"`
}

var Config *config

func initialiseViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	viper.AddConfigPath("/etc/ekko")
	if err := viper.ReadInConfig(); err != nil {
		log.Panicf("Error reading config file: %s\n", err)
	}
}

func init() {
	initialiseViper()
	if err := viper.Unmarshal(&Config); err != nil {
		log.Panicf("Invalid configuration, %s", err)
	}
	defaults.SetDefaults(Config)
	if viper.GetBool("logging.console_enabled") && viper.GetBool("ui_enabled") {
		log.Panicf("Both logging.console_enabled & ui_enabled can't be enabled simultaneously")
	}
}
