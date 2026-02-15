package main

import (
	"log"

	"github.com/ftp27/github-readme-stats/internal/config"
	"github.com/ftp27/github-readme-stats/internal/server"
	"github.com/spf13/viper"
)

func main() {
	config.Load()

	router := server.NewRouter()

	port := viper.GetString("PORT")
	if port == "" {
		port = "9000"
	}

	addr := "0.0.0.0:" + port
	if err := router.Run(addr); err != nil {
		log.Fatal(err)
	}
}
