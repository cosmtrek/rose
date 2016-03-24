package main

import (
	"flag"
	"os"
	"path"
	"strconv"

	"github.com/BurntSushi/toml"
)

type Config struct {
	ServerEnv     string `toml:"server_env"`
	ServerHost    string `toml:"server_host"`
	ServerPort    string `toml:"server_port"`
	SocketTimeout int    `toml:"socket_timeout"`
	PidFile       string `toml:"pid_file,omitempty"`
}

const (
	RoseToml = ".rose.toml"
)

var (
	config = Config{
		ServerEnv:     "development",
		ServerHost:    "127.0.0.1",
		ServerPort:    "3333",
		SocketTimeout: 300,
		PidFile:       "./.rose.pid",
	}
)

func initConfig() {
	loadConfigFile()
	overrideConfig()
	initLog()
	showConfig()
}

func loadConfigFile() {
	// Default config file is ~/.rose.toml
	file := getDefaultRoseToml()
	if _, err := os.Stat(file); err == nil {
		if _, err := toml.DecodeFile(file, &config); err != nil {
			errl.Printf("Failed to parse ~/.rose.toml, error: %s", err)
		}
	}
}

func overrideConfig() {
	flag.StringVar(&config.ServerEnv, "server_env", config.ServerEnv, "server environment")
	flag.StringVar(&config.ServerHost, "server_host", config.ServerHost, "server host")
	flag.StringVar(&config.ServerPort, "server_port", config.ServerPort, "server port")
	flag.IntVar(&config.SocketTimeout, "socket_timeout", config.SocketTimeout, "socket time out")
	flag.Parse()
}

func showConfig() {
	info.Printf("Server %s %s:%s is starting...", config.ServerEnv, config.ServerHost, config.ServerPort)
	info.Println("Socket time out: " + strconv.Itoa(config.SocketTimeout))
}

func getDefaultRoseToml() string {
	return path.Join(os.Getenv("HOME"), RoseToml)
}
