package main

import (
	"GB_Study_02/fin/base/sqlite"
	"GB_Study_02/fin/conf"
	"GB_Study_02/fin/web"
	"context"
	"log"
)

var LastId int

type Cfg struct {
	Address  string
	Port     string
	DBFile   string
	JSFile   string
	Template string
}

func main() {

	var cf conf.Config

	Cfg := (*conf.Config).Conf(&cf)

	//	Cfg := *c

	s, err := sqlite.New(Cfg.DBFile)

	if err != nil {
		log.Fatalf("can't connect to DB: %s", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("can't init DB: %s", err)
	}

	if h, err := web.NewHttp(Cfg.Address, Cfg.Port); err != nil {
		panic(err)
	} else {
		log.Println(h)
	}

}
