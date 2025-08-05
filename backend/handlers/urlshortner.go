package handlers

import (
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func UrlShortenHandler(pool *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		var req struct {
			URL string `json:"url"`
		}
		if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
			logger.Error("Failed to decode request", "error", err)
			http.Error(w, "Invalid request body", http.StatusBadRequest)
			return
		}

		if req.URL == "" {
			logger.Warn("Empty URL")
			http.Error(w, "URL is required", http.StatusBadRequest)
			return
		}

		userID := interface{}(nil)
		if token := r.Header.Get("Authorization"); token != "" {
			if t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
				return []byte(os.Getenv("JWT_SECRET")), nil
			}); err == nil && t.Valid {
				if claims, ok := t.Claims.(jwt.MapClaims); ok {
					userID = int(claims["user_id"].(float64))
				}
			}
		}

		shortKey := uuid.New().String()[:8]
		ctx := context.Background()
		_, err := pool.Exec(ctx, `
			INSERT INTO urls (short_key, original_url, user_id, created_at, expires_at)
			VALUES ($1, $2, $3, $4, $5)`,
			shortKey, req.URL, userID, time.Now(), time.Now().Add(24*time.Hour))
		if err != nil {
			logger.Error("Failed to insert URL", "error", err)
			http.Error(w, "Failed to shorten URL", http.StatusInternalServerError)
			return
		}

		resp := struct {
			ShortURL string `json:"shortUrl"`
		}{ShortURL: "http://localhost:8080/url/" + shortKey}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(resp); err != nil {
			logger.Error("Failed to encode response", "error", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
		}
	}
}

func UrlRedirectHandler(pool *pgxpool.Pool, logger *slog.Logger) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
			return
		}

		shortKey := strings.TrimPrefix(r.URL.Path, "/url/")
		if shortKey == "" {
			logger.Warn("Empty short key")
			http.Error(w, "Short key is required", http.StatusBadRequest)
			return
		}

		var originalURL string
		ctx := context.Background()
		err := pool.QueryRow(ctx, `
			SELECT original_url FROM urls
			WHERE short_key = $1 AND expires_at > $2`,
			shortKey, time.Now()).Scan(&originalURL)
		if err == pgx.ErrNoRows {
			logger.Warn("URL not found or expired", "short_key", shortKey)
			http.Error(w, "URL not found or expired", http.StatusNotFound)
			return
		} else if err != nil {
			logger.Error("Failed to query URL", "error", err, "short_key", shortKey)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, originalURL, http.StatusMovedPermanently)
	}
}