package main

import (
	_ "server/docs"
	"server/internal/adapters/adrepo"
	"server/internal/app"
	"server/internal/ports/httpgin"
)

// @title           Ad server documentation
// @version         1.0
// @host      localhost:18080
// @BasePath  /
func main() {
	server := httpgin.NewHTTPServer(":18080", app.NewApp(adrepo.New()))
	err := server.Listen()
	if err != nil {
		panic(err)
	}
}
