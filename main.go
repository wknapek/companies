package main

import (
	"flag"
	"net/http"
	"os"
	"time"

	"companies/config"
	"companies/handler"
	"companies/security"
	"github.com/rs/zerolog/log"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func main() {
	log.Logger = zerolog.New(zerolog.ConsoleWriter{Out: os.Stderr, TimeFormat: time.RFC3339}).
		Level(zerolog.TraceLevel).
		With().
		Timestamp().
		Caller().
		Int("pid", os.Getpid()).
		Logger()
	pathToConf := flag.String("conf", config.ConfFile, "path to configuration")
	var help = flag.Bool("help", false, "Show help")
	flag.Parse()
	if *help {
		flag.Usage()
		os.Exit(0) //nolint:gocritic
	}
	cfg, errConf := config.ReadConfiguration(*pathToConf)
	if errConf != nil {
		log.Error().Msgf("no configuration file %s", *pathToConf)
		os.Exit(0)
	}
	log.Info().Msg("app runing")
	sec := security.New(map[string]string{"admin": "admin"}, cfg.SessionTime, []byte(cfg.JWT_Key))
	router := chi.NewRouter()
	handl := handler.NewHandler(cfg.DBURI, cfg.DBUser, cfg.DBPasswd, sec)
	handl.Init()
	router.Use(middleware.Logger)
	router.Post("/companies/login", sec.Login)
	router.Delete("/companies/{name}", handl.Delete)
	router.Patch("/companies", handl.Update)
	router.Post("/companies", handl.Create)
	router.Get("/companies/{name}", handl.Get)
	errLAndS := http.ListenAndServe(":"+cfg.Port, router)
	if errLAndS != nil {
		return
	}
}
