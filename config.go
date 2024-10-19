package main

import (
	"encoding/json"
	"fmt"
	"os"
)

type ServerConfig struct {
	Port     int    `json:"port"`
	BindAddr string `json:"bind_address"`
}

type Config struct {
	MarkdownDir string       `json:"markdown_dir"`
	AssetsDir   string       `json:"assets_dir"`
	Server      ServerConfig `json:"server"`
}

func (c *Config) read(cfg_path string) []byte {
	config_data, err := os.ReadFile(cfg_path)
	if err != nil {
		panic(fmt.Sprintf("Unable to read the config file. Error: %s", err))
	}

	return config_data
}

func (c *Config) Init() {
	cfg_path := os.Getenv("GOBLOG_CONFIG_PATH")
	if cfg_path == "" {
		cfg_path = "./.goblog.json"
	}
	config_data := c.read(cfg_path)

	err := json.Unmarshal(config_data, c)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse config file. Error: %s", err))
	}
}
