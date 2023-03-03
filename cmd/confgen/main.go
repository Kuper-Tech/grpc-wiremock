package main

import (
	"log"

	"github.com/SberMarket-Tech/grpc-wiremock/cmd/confgen/commands"
)

func main() {
	command := commands.CreateCommandRoot()

	if err := command.Execute(); err != nil {
		log.Fatalln("execute cli command:", err.Error())
	}
}
