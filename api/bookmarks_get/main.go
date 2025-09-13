package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/redis/go-redis/v9"
)

var (
	ctx3 = context.Background()
	rdb3 = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
)

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusBadRequest)
		return
	}

	username, err := rdb3.Get(ctx3, "session:"+token).Result()
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	key := fmt.Sprintf("bookmarks:%s", username)
	bookmarks, err := rdb3.LRange(ctx3, key, 0, -1).Result()
	if err != nil {
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(bookmarks)
}
