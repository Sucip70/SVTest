package main

import (
	"log"
	"net/http"

	"github.com/Sucip70/SVTest-demo-be/internal/config"
	"github.com/Sucip70/SVTest-demo-be/internal/db"
	"github.com/Sucip70/SVTest-demo-be/internal/handler"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Println("No .env file found, using environment variables")
	}

	cfg := config.Load()
	database, err := db.Connect(cfg.DBDSN)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer database.Close()

	err = db.RunMigrations(database)
	if err != nil {
		log.Fatalf("Failed to run migrations: %v", err)
	}

	postsHandler := &handler.PostsHandler{DB: database}

	mux := http.NewServeMux()
	mux.HandleFunc("GET /posts", postsHandler.GetPosts)
	mux.HandleFunc("GET /post", postsHandler.GetPostByID)
	mux.HandleFunc("POST /post", postsHandler.CreatePost)
	mux.HandleFunc("PUT /post", postsHandler.UpdatePost)
	mux.HandleFunc("PATCH /post/status", postsHandler.UpdatePostStatus)

	mux.HandleFunc("GET /health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	handler := corsMiddleware(mux)

	port := cfg.Port
	if port == "" {
		port = "8080"
	}

	log.Printf("Server is running on :%s", port)
	if err := http.ListenAndServe(":"+port, handler); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Header.Get("Origin") == "https://sucip70.github.io" {
			w.Header().Set("Access-Control-Allow-Origin", "https://sucip70.github.io")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, PATCH, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
			w.Header().Set("Vary", "Origin")
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}

