package main

import (
	"flag"
	"github.com/ogra1/fabrica/config"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/datastore/sqlite"
	"github.com/ogra1/fabrica/service/repo"
	"github.com/ogra1/fabrica/service/watch"
	"github.com/ogra1/fabrica/web"
)

func main() {
	var mode string
	flag.StringVar(&mode, "mode", "web", "Mode of operation: web, watch")
	flag.Parse()

	settings := config.ReadParameters()

	// Set up the dependency chain
	db, _ := sqlite.NewDatabase()
	buildSrv := repo.NewBuildService(db)

	// Set up the service based on the mode
	if mode == "watch" {
		watchDaemon(db, buildSrv)
	} else {
		webService(settings, buildSrv)
	}
}

func webService(settings *config.Settings, buildSrv *repo.BuildService) {
	srv := web.NewWebService(settings, buildSrv)
	srv.Start()
}

func watchDaemon(db datastore.Datastore, buildSrv repo.BuildSrv) {
	watchSrv := watch.NewWatchService(db, buildSrv)
	watchSrv.Watch()
}
