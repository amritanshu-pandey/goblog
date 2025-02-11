package main

import "amritanshu.in/goblog/backend"

func main() {
	config := Config{}
	config.Init()

	backend.RunServer(config.MarkdownDir, config.AssetsDir, config.Server.Port, config.Server.BindAddr, config.ShowDrafts)
}
