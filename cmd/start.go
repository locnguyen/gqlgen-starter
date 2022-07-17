package main

import (
	"encoding/json"
	"gqlgen-starter/config"
	"gqlgen-starter/internal"
	"log"
)

func main() {
	log.Printf("******************************************")
	log.Printf("\tBuild Commit: %s", BuildCommit)
	log.Printf("\tBuild Time: %s", BuildTime)
	configJson, err := json.MarshalIndent(config.Application, "", "  ")
	if err != nil {
		log.Fatalf(err.Error())
	}
	log.Printf("\tEnvironment Variables: %s", string(configJson))
	log.Printf("******************************************")

	internal.StartServer()
}
