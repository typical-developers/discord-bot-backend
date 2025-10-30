package main

import (
	"context"
	"database/sql"
	"fmt"
	"net/http"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-chi/chi"
	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	_ "github.com/lib/pq"
	"github.com/redis/go-redis/v9"
	"github.com/sirupsen/logrus"
	log "github.com/sirupsen/logrus"
	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/typical-developers/discord-bot-backend/internal/db"
	_ "github.com/typical-developers/discord-bot-backend/internal/logger"
	discord_state "github.com/typical-developers/discord-bot-backend/pkg/discord-state"
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

func discordRedisConnect() (*redis.Client, error) {
	var opts = &redis.Options{
		Addr:     fmt.Sprintf("%s:%d", config.C.DiscordCache.Host, config.C.DiscordCache.Port),
		Password: config.C.DiscordCache.Password,
		DB:       config.C.DiscordCache.DB,
	}

	client := redis.NewClient(opts)

	ticker := time.NewTicker(time.Second * 10)
	go func() {
		ctx := context.Background()
		defer ticker.Stop()

		healthy := true
		for range ticker.C {
			_, err := client.Ping(ctx).Result()

			if err != nil {
				healthy = false
				log.Warn("Redis client connection is not healthy.")
				continue
			}

			if !healthy {
				healthy = true
				log.Info("Redis client connection has been restored.")
				continue
			}

			log.Debug("Redis client connection is healthy.")
		}
	}()

	return client, nil
}

func serveStatic(r *chi.Mux) {
	assetsRoot := http.Dir("./assets")
	fs := http.StripPrefix("/static/", http.FileServer(assetsRoot))

	r.Handle("/static/*", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "*")

		fs.ServeHTTP(w, r)
	}))
}

func runMigrations(db *sql.DB) error {
	driver, err := postgres.WithInstance(db, &postgres.Config{})
	if err != nil {
		return err
	}

	m, err := migrate.NewWithDatabaseInstance(
		"file://internal/db/migrations",
		"postgres", driver,
	)

	if err != nil {
		return err
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}

	if err == migrate.ErrNoChange {
		log.Info("No database migrations to run.")
	} else {
		log.Info("Successfully ran database migrations.")
	}

	return nil
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
	err = runMigrations(pqdb)
	if err != nil {
		panic(err)
	}

	querier := db.New(pqdb)

	discord, err := discordgo.New("Bot " + config.C.DiscordToken)
	if err != nil {
		panic(err)
	}

	discord.Identify.Intents = discordgo.IntentsGuilds |
		discordgo.IntentsGuildMembers |
		discordgo.IntentsGuildMessages
	err = discord.Open()
	if err != nil {
		panic(err)
	}

	discordCache, err := discordRedisConnect()
	if err != nil {
		panic(err)
	}
	discordState := discord_state.NewStateManager(&discord_state.StateManagerOptions{
		DiscordSession: discord,
		RedisClient:    discordCache,
	})

	guildUsecase := usecase.NewGuildUsecase(pqdb, querier, discordState)
	handlers.NewGuildHandler(router, guildUsecase)

	memberUsecase := usecase.NewMemberUsecase(pqdb, querier, discordState)
	handlers.NewMemberHandler(router, memberUsecase)

	port := fmt.Sprintf(":%d", config.C.Port)
	panic(http.ListenAndServe(port, router))
}
