package main

import (
	"github.com/mcuadros/go-defaults"
	"github.com/soheltarir/gollections/containers"
	"github.com/spf13/viper"
	"log"
)

// Server defines the object of a Game server
type Server struct {
	Name      string `mapstructure:"name"`
	Game      string `mapstructure:"game"`
	IPAddress string `mapstructure:"ip"`
}

func (s Server) Key() interface{} {
	return s.Name
}

func (s Server) Less(value containers.Container) bool {
	y := value.(Server)
	return s.Name < y.Name
}

func (Server) Validate(x interface{}) containers.Container {
	return x.(Server)
}

type loggingConfig struct {
	FileEnabled    bool   `mapstructure:"file_enabled"`
	ConsoleEnabled bool   `mapstructure:"console_enabled"`
	FileOutput     string `mapstructure:"file_output"`
}

type config struct {
	Servers        []Server `mapstructure:"servers"`
	Logging        loggingConfig
	MaxPacketNum   int   `mapstructure:"max_packet_num" default:"20"`
	MinPacketNum   int   `mapstructure:"min_packet_num" default:"4"`
	PingTimeout    int64 `mapstructure:"ping_timeout" default:"30"` // in seconds
	WorkerPoolSize int   `mapstructure:"worker_pool_size" default:"5"`
}

var Config *config

func initialiseViper() {
	viper.SetConfigName("config")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
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
	if !viper.IsSet("logging.file_output") {
		Config.Logging.FileOutput = DefaultFileLogPath
	}
}
