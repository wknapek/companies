package main

import (
	"net/http"

	"companies/handler"
	"companies/security"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/rs/zerolog"
)

func main() {
	zerolog.SetGlobalLevel(zerolog.InfoLevel)
	sec := security.Security{}
	router := chi.NewRouter()
	handl := handler.NewHandler("", "", "")
	router.Use(middleware.Logger)
	router.Post("/companies/login", sec.Login)
	router.Delete("/companies/{name}", handl.Delete)
	router.Patch("/companies", handl.Update)
	router.Post("/companies", handl.Create)
	router.Get("/companies/{name}", handl.Get)
	err := http.ListenAndServe(":3001", router)
	if err != nil {
		return
	}
}
