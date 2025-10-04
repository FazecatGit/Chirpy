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
	JWTSecret      string
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
	apiCfg.JWTSecret = os.Getenv("JWT_SECRET")

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
	mux.HandleFunc("/api/login", methodHandler("POST", apiCfg.handlerLoginUser))
	mux.HandleFunc("/api/refresh", methodHandler("POST", apiCfg.handlerRefreshToken))
	mux.HandleFunc("/api/revoke", methodHandler("POST", apiCfg.handlerRevokeToken))
	mux.HandleFunc("/api/chirps/{chirp_id}", methodHandler("DELETE", apiCfg.deleteChirpHandler))

	mux.HandleFunc("/api/chirps", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "POST":
			apiCfg.createChirpHandler(w, r)
		case "GET":
			apiCfg.listallChirpsHandler(w, r)
		default:
			w.Header().Set("Allow", "GET, POST")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/users", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			apiCfg.handlerCreateUser(w, r)
		case http.MethodPut:
			apiCfg.handlerAuthorizeUser(w, r)
		default:
			w.Header().Set("Allow", "POST, PUT")
			http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
		}
	})

	mux.HandleFunc("/api/chirps/{chirp_id}", methodHandler("GET", apiCfg.getChirpHandler))

	server := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
