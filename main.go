package main

import (
	"database/sql"
	"github.com/go-chi/chi"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/peyzor/rssagg/internal/database"
	"log"
	"net/http"
	"os"
	"time"
)

type serverConfig struct {
	DB *database.Queries
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL is not found")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("can't connect to database")
	}

	db := database.New(conn)
	serverConfig := serverConfig{
		DB: db,
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT is not found")
	}

	go startScraping(db, 10, time.Minute)

	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()
	v1Router.Get("/healthz", handlerReadiness)
	v1Router.Get("/err", handlerErr)
	v1Router.Post("/users", serverConfig.handlerCreateUser)
	v1Router.Get("/users", serverConfig.middlewareAuth(serverConfig.handlerGetUser))
	v1Router.Post("/feeds", serverConfig.middlewareAuth(serverConfig.handlerCreateFeed))
	v1Router.Get("/feeds", serverConfig.handlerGetFeeds)
	v1Router.Post("/feeds/follow", serverConfig.middlewareAuth(serverConfig.handlerCreateFeedFollow))
	v1Router.Get("/feeds/follow", serverConfig.middlewareAuth(serverConfig.handlerGetFeedFollows))
	v1Router.Delete("/feeds/follow/{feedFollowID}", serverConfig.middlewareAuth(serverConfig.handlerDeleteFeedFollow))

	router.Mount("/v1", v1Router)

	srv := &http.Server{
		Handler: router,
		Addr:    ":" + port,
	}

	log.Printf("Server starting on port %s", port)
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}
