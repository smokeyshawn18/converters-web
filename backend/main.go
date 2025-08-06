package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"syscall"
	"time"

	"utility-app/db"
	"utility-app/handlers"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
)

func corsMiddleware(next http.Handler, allowedOrigin string) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        origin := r.Header.Get("Origin")
        allowedOrigins := []string{
            "http://localhost:3000",
			"https://converters-web.vercel.app/",
            "http://127.0.0.1:3000",
            "http://frontend:3000", // For Docker
        }
        for _, allowed := range allowedOrigins {
            if origin == allowed {
                w.Header().Set("Access-Control-Allow-Origin", origin)
                break
            }
        }
        w.Header().Set("Access-Control-Allow-Methods", "GET, POST, OPTIONS")
        w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Accept, Authorization, authorization")
        w.Header().Set("Access-Control-Max-Age", "86400")
        if r.Method == http.MethodOptions {
            w.WriteHeader(http.StatusOK)
            return
        }
        next.ServeHTTP(w, r)
    })
}

func rateLimitMiddleware(next http.Handler, requests int, window time.Duration) http.Handler {
	limiter := rate.NewLimiter(rate.Every(window), requests)
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := limiter.Wait(context.Background()); err != nil {
			http.Error(w, "Rate limit exceeded", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}))
	slog.SetDefault(logger)

if os.Getenv("RENDER") == "" {
	if err := godotenv.Load(); err != nil {
		logger.Warn("No .env file found (running locally?)", "error", err)
	}
}


	allowedOrigin := os.Getenv("ALLOWED_ORIGIN")
	if allowedOrigin == "" {
		allowedOrigin = "http://localhost:3000"
		logger.Warn("ALLOWED_ORIGIN not set, defaulting to http://localhost:3000")
	}

	dbURL := os.Getenv("DATABASE_URL")
	if dbURL == "" {
		logger.Error("DATABASE_URL not set")
		os.Exit(1)
	}
	pool, err := pgxpool.New(context.Background(), dbURL)
	if err != nil {
		logger.Error("Failed to connect to database", "error", err)
		os.Exit(1)
	}
	defer pool.Close()

	if err := db.InitializeDatabase(pool, logger); err != nil {
		logger.Error("Failed to initialize database", "error", err)
		os.Exit(1)
	}

	rateLimitRequests, _ := strconv.Atoi(os.Getenv("RATE_LIMIT_REQUESTS"))
	if rateLimitRequests == 0 {
		rateLimitRequests = 100
	}
	rateLimitWindow, _ := time.ParseDuration(os.Getenv("RATE_LIMIT_WINDOW"))
	if rateLimitWindow == 0 {
		rateLimitWindow = time.Minute
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/auth/register", handlers.RegisterHandler(pool, logger))
	mux.HandleFunc("/auth/login", handlers.LoginHandler(pool, logger))
mux.HandleFunc("/youtube/mp3", handlers.YoutubeToMP3Handler(logger))

	mux.HandleFunc("/url/shorten", handlers.UrlShortenHandler(pool, logger))
	mux.HandleFunc("/url/", handlers.UrlRedirectHandler(pool, logger))
	mux.HandleFunc("/image/convert", handlers.ImageConvertHandler(pool, logger))
	mux.HandleFunc("/document/convert", handlers.DocConvertHandler(pool, logger))


	handler := corsMiddleware(mux, allowedOrigin)
	handler = rateLimitMiddleware(handler, rateLimitRequests, rateLimitWindow)

	server := &http.Server{
		Addr:         ":8080",
		Handler:      handler,
		ReadTimeout:  30 * time.Second,
		WriteTimeout: 5 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	go func() {
		sigChan := make(chan os.Signal, 1)
		signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
		<-sigChan
		logger.Info("Received shutdown signal, stopping server...")
		ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
		defer cancel()
		if err := server.Shutdown(ctx); err != nil {
			logger.Error("Server shutdown failed", "error", err)
			os.Exit(1)
		}
	}()

	logger.Info("Starting server on :8080")
	if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("Server failed to start", "error", err)
		os.Exit(1)
	}
}