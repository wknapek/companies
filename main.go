package main

import (
	"flag"
	"net/http"
	"os"

	"companies/config"
	"companies/handler"
	"companies/security"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func main() {
	pathToConf := flag.String("conf", config.ConfFile, "path to configuration")
	var help = flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0) //nolint:gocritic
	}
	cfg, errConf := config.ReadConfiguration(*pathToConf)
	if errConf != nil {
		panic(errConf)
	}
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	sec := security.New(map[string]string{"admin": "admin"}, cfg.SessionTime)
	router := chi.NewRouter()
	handl := handler.NewHandler(cfg.DBURI, cfg.DBUser, cfg.DBPasswd)
	handl.Init()
	router.Use(middleware.Logger)
	router.Post("/companies/login", sec.Login)
	router.Delete("/companies/{name}", handl.Delete)
	router.Patch("/companies", handl.Update)
	router.Post("/companies", handl.Create)
	router.Get("/companies/{name}", handl.Get)
	errLAndS := http.ListenAndServe(":3001", router)
	if errLAndS != nil {
		return
	}
}
