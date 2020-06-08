package main

import (
	"flag"
	"github.com/ogra1/fabrica/config"
	"github.com/ogra1/fabrica/datastore"
	"github.com/ogra1/fabrica/datastore/sqlite"
	"github.com/ogra1/fabrica/service/key"
	"github.com/ogra1/fabrica/service/lxd"
	"github.com/ogra1/fabrica/service/repo"
	"github.com/ogra1/fabrica/service/system"
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
	systemSrv := system.NewSystemService(db)
	lxdSrv := lxd.NewLXD(db, systemSrv)
	buildSrv := repo.NewBuildService(db, lxdSrv)
	keySrv := key.NewKeyService(db)

	// Set up the service based on the mode
	if mode == "watch" {
		watchDaemon(db, buildSrv, keySrv)
	} else {
		webService(settings, buildSrv, lxdSrv, systemSrv, keySrv)
	}
}

func webService(settings *config.Settings, buildSrv *repo.BuildService, lxdSrv lxd.Service, systemSrv system.Srv, keySrv key.Srv) {
	srv := web.NewWebService(settings, buildSrv, lxdSrv, systemSrv, keySrv)
	srv.Start()
}

func watchDaemon(db datastore.Datastore, buildSrv repo.BuildSrv, keySrv key.Srv) {
	watchSrv := watch.NewWatchService(db, buildSrv, keySrv)
	watchSrv.Watch()
}
