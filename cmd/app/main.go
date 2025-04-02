//	@title			ChillGuys API
//	@version		1.0
//	@description	API for ChillGuys marketplace
//	@host			90.156.217.63:8081
//	@BasePath		/api

//	@securityDefinitions.basic	BasicAuth
//	@securityDefinitions.apikey	TokenAuth
//	@in							cookie
//	@name						token

package main

import (
	_ "github.com/go-park-mail-ru/2025_1_ChillGuys/docs"
	"github.com/go-park-mail-ru/2025_1_ChillGuys/internal/app"
	_ "github.com/lib/pq"
	"log"
)

func main() {
	if err := app.Run(); err != nil {
		log.Fatalf("Application failed: %v", err)
	}
}
