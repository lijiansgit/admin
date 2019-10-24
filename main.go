package main

import (
	"flag"

	"github.com/lijiansgit/admin/pkg/ldap"

	"github.com/lijiansgit/admin/config"
	"github.com/lijiansgit/admin/models"
	"github.com/lijiansgit/admin/routers"
	log "github.com/lijiansgit/go/libs/log4go"
)

func main() {
	flag.Parse()

	if err := config.Init(); err != nil {
		panic(err)
	}

	log.LoadConfiguration(config.Conf.Log.Conf)
	defer log.Close()

	log.Debug("Load config file %s: %#v", config.ConfFile, config.Conf)

	if err := ldap.Init(); err != nil {
		panic(err)
	}

	if err := models.Init(); err != nil {
		panic(err)
	}

	r := routers.GetRouters()
	r.Run(config.Conf.WEB.Addr)
}
