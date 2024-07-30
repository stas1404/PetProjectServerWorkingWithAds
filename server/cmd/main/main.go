package main

import (
	"server/internal/adapters/adrepo"
	"server/internal/app"
	"server/internal/ports/httpgin"
)

func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
