package http

import(
	"context"
	"encoding/json"
	"net/http"
	"time"

	"github.com/go-redis/redis/v8"
)

type API struct {
	Redis *redis.Client
}

func (a *API) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("POST /api/v1/shorten", a.ShortenHandler)
	mux.HandleFunc("GET /{key}", a.RedirectHandler)
	return mux
}


func (a *API) ShortenHandler(w http.ResponseWriter, r *http.Request) {
	var payload struct {
		URL string `json:"url"`
	}

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if payload.URL == "" {
		http.Error(w, "URL cannot be empty", http.StatusBadRequest)
		return
	}

	normalized := Normalize(payload.URL)
	hash := ShortHash(normalized)

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	finalHash := hash
	maxCollisionRetries := 10
	for attempts := 0; attempts < maxCollisionRetries; attempts++ { 
		existing, err := a.Redis.Get(ctx, finalHash).Result()
		if err == redis.Nil {
			break
		}

		if err != nil {
			http.Error(w, "Redis GET error", http.StatusInternalServerError)
			return
		}

		if existing == normalized {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(map[string]string{
				"key":       finalHash,
				"long_url":  normalized,
				"short_url": getShortURLBase() + "/" + finalHash,
			})
			return
		}
		finalHash = finalHash + generateRandomSuffix(attempts + 1)
	}

	if err := a.Redis.Set(ctx, finalHash, normalized, 0).Err(); err != nil {
		http.Error(w, "Redis SET error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"key":       finalHash,
		"long_url":  normalized,
		"short_url": getShortURLBase() + "/" + finalHash,
	})
}

func (a *API) RedirectHandler(w http.ResponseWriter, r *http.Request) {
	key := r.PathValue("key")
	if key == "" {
		http.Error(w, "Missing key in URL path", http.StatusBadRequest)
		return
	}

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	val, err := a.Redis.Get(ctx, key).Result()
	if err == redis.Nil {
		http.NotFound(w, r)
		return
	} else if err != nil {
		http.Error(w, "Failed to get value from Redis", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, val, http.StatusFound)
}
