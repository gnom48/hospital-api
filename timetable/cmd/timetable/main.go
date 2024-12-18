package main

import (
	"flag"
	"log"
	"timetable/internal/server"

	"github.com/BurntSushi/toml"
)

/*
Запуск
go run cmd/account/main.go

Сборка документации
cd internal/server
swag init -g server.go appointment.go timetable.go
*/

var configPath string

func init() {
	flag.StringVar(&configPath, "config-path", "config/config.toml", "Path to configure file")
}

func main() {
	flag.Parse()

	config := server.NewConfig()
	_, e := toml.DecodeFile(configPath, config)
	if e != nil {
		log.Fatal(e)
	}

	server := server.New(config)
	if err := server.Start(); err == nil {
		log.Fatal(err)
	}
}
