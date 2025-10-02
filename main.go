package main

//go build -o out && ./out

import (
	"database/sql"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"

	"github.com/FazecatGit/Chirpy/internal/database"
	"github.com/google/uuid"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

type apiConfig struct {
	fileServerHits atomic.Int32
	DB             *database.Queries
	Platform       string
}

type User struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Email     string    `json:"email"`
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found")
	}

	dbURL := os.Getenv("DB_URL")
	db, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Cannot connect to database:", err)
	}

	apiCfg := &apiConfig{
		DB:       database.New(db),
		Platform: os.Getenv("PLATFORM"),
	}

	fileServer := http.FileServer(http.Dir("."))
	mux := http.NewServeMux()

	mux.HandleFunc("/api/healthz", methodHandler("GET", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	}))

	mux.Handle("/app/", apiCfg.middlewareMetricsInc(http.StripPrefix("/app", fileServer)))
	mux.HandleFunc("/admin/metrics", methodHandler("GET", apiCfg.metricsHandler))
	mux.HandleFunc("/admin/reset", methodHandler("POST", apiCfg.adminResetHandler))
	mux.HandleFunc("/api/validate_chirp", methodHandler("POST", apiCfg.validateChirpHandler))
	mux.HandleFunc("/api/users", methodHandler("POST", apiCfg.handlerCreateUser))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
