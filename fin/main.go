package main

import (
	"GB_Study_02/fin/base/sqlite"
	"GB_Study_02/fin/conf"
	"GB_Study_02/fin/web"
	"context"
	"log"
)

var LastId int

func main() {

	var cf conf.Config

	c := (*conf.Config).Conf(&cf)

	s, err := sqlite.New(c.DBFile)

	if err != nil {
		log.Fatalf("can't connect to DB: %s", err)
	}

	if err := s.Init(context.TODO()); err != nil {
		log.Fatalf("can't init DB: %s", err)
	}

	if h, err := web.NewHttp(c.Address, c.Port); err != nil {
		panic(err)
	} else {
		log.Println(h)
	}

}
