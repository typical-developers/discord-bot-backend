package main

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/go-chi/chi"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	_ "github.com/typical-developers/discord-bot-backend/internal/logger"
	"github.com/typical-developers/discord-bot-backend/services/web/config"
	_ "github.com/typical-developers/discord-bot-backend/services/web/config"
	_ "github.com/typical-developers/discord-bot-backend/services/web/docs"
	"github.com/typical-developers/discord-bot-backend/services/web/handlers"
	"github.com/typical-developers/discord-bot-backend/services/web/usecase"
)

func dbConnect() (*sql.DB, error) {
	db, err := sql.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:%d?%s",
		config.C.Database.Username,
		config.C.Database.Password,
		config.C.Database.Host,
		config.C.Database.Port,
		config.C.Database.Options,
	))

	if err != nil {
		return nil, err
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(50)

	return db, nil
}

func serveStatic(r *chi.Mux) {
	assetsRoot := http.Dir("./assets")
	fs := http.StripPrefix("/static/", http.FileServer(assetsRoot))

	r.Handle("/static/*", fs)
}

//	@title						Discord Bot API
//	@version					1.0
//	@description				The API for the main Typical Developers Discord bot.
//
//	@tag.name					Guilds
//	@tag.description			Guild endpoints.
//
//	@tag.name					Voice Room Lobbies
//	@tag.description			Voice room lobby endpoints.
//
//	@tag.name					Voice Rooms
//	@tag.description			Voice room endpoints.
//
//	@tag.name					Members
//	@tag.description			Member endpoints.
//
//	@tag.name					HTML Generation
//	@tag.description			HTML generation endpoints.
//
//	@securitydefinitions.apikey	APIKeyAuth
//	@in							header
//	@name						X-API-KEY
//
// nolint:staticcheck
func main() {
	if lvl, err := logrus.ParseLevel(config.C.LogLevel); err != nil {
		logrus.SetLevel(lvl)
	}

	router := chi.NewRouter()
	router.Use(handlers.RequestLog)
	router.Get("/docs/*", httpSwagger.Handler())
	serveStatic(router)

	pqdb, err := dbConnect()
	if err != nil {
		panic(err)
	}
	querier := db.New(pqdb)

	guildUsecase := usecase.NewGuildUsecase(querier)
	handlers.NewGuildHandler(router, guildUsecase)

	port := fmt.Sprintf(":%d", config.C.Port)
	panic(http.ListenAndServe(port, router))
}
