package main

import (
	"database/sql"
	"fmt"
	config "go-bank-app/configs"
	"go-bank-app/internal/auth"
	"go-bank-app/pkg/middleware"
	"log"
	"net/http"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("failed to load .env file")
	}

	dbURL, err := config.GetString("POSTGRES_DB_URI")
	if err != nil {
		log.Fatal(err)
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("failed to connect to DB:", err)
	}
	defer conn.Close()

	authRepo := auth.NewAuthRepository(conn)
	authService := auth.NewAuthService(authRepo)
	authHandler := auth.NewAuthHandler(authService)

	// Rutas públicas
	http.HandleFunc("/auth/register", authHandler.Register)
	http.HandleFunc("/auth/login", authHandler.Login)

	// Ruta protegida con middleware
	http.Handle("/auth/me", middleware.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		userID := r.Context().Value(middleware.ContextUserIDKey).(string)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		fmt.Fprintf(w, `{"user_id": "%s"}`, userID)
	})))

	port := ":8070"
	fmt.Println("🚀 Server running at http://localhost" + port)
	log.Fatal(http.ListenAndServe(port, nil))
}
