//	@title			ChillGuys API
//	@version		1.0
//	@description	API for ChillGuys marketplace
//	@host			90.156.217.63:8081
//	@BasePath		/api/v1

//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	TokenAuth
//	@in							cookie
//	@name						token

package main

import (
	"log"

	"github.com/go-park-mail-ru/2025_1_ChillGuys/config"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/app"
	_ "github.com/lib/pq"
)

func main() {
	conf, err := config.NewConfig()
	if err != nil {
		log.Fatalf("config error: %v", err)
	}

	application, err := app.NewApp(conf)
	if err != nil {
		log.Fatalf("failed to create app: %v", err)
	}

	if err := application.Run(); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
