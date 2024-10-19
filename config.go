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

func (c *Config) read(cfgPath string) []byte {
	configData, err := os.ReadFile(cfgPath)
	if err != nil {
		panic(fmt.Sprintf("Unable to read the config file. Error: %s", err))
	}

	return configData
}

func (c *Config) Init() {
	cfgPath := os.Getenv("GOBLOG_CONFIG_PATH")
	if cfgPath == "" {
		cfgPath = "./.goblog.json"
	}
	configData := c.read(cfgPath)

	err := json.Unmarshal(configData, c)
	if err != nil {
		panic(fmt.Sprintf("Unable to parse config file. Error: %s", err))
	}
}
