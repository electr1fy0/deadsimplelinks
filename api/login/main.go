package handler

import (
	"context"
	"encoding/json"
	"os"

	"net/http"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

var (
	ctx = context.Background()

	rdb = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
)

type LoginRequest struct {
	Username string `json:"username"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest

	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	token := uuid.NewString()
	if err := rdb.Set(ctx, "session:"+token, req.Username, 0).Err(); err != nil {
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{"token": token})
}
