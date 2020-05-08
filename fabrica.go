package main

import (
	"github.com/ogra1/fabrica/config"
	"github.com/ogra1/fabrica/datastore/sqlite"
	"github.com/ogra1/fabrica/service"
	"github.com/ogra1/fabrica/web"
)

func main() {
	settings := config.ReadParameters()

	// Set up the dependency chain
	db, _ := sqlite.NewDatabase()
	buildSrv := service.NewBuildService(db)

	srv := web.NewWebService(settings, buildSrv)
	srv.Start()
}
