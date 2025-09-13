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
	ctx2 = context.Background()
	rdb2 = redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_ADDR"),
		Username: "default",
		Password: os.Getenv("REDIS_PASS"),
		DB:       0,
	})
)

type BookmarkRequest struct {
	Token string `json:"token"`
	URL   string `json:"url"`
}

func Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "invalid method", http.StatusMethodNotAllowed)
		return
	}
	var req BookmarkRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request", http.StatusBadRequest)
		return
	}

	username, err := rdb2.Get(ctx2, "session:"+req.Token).Result()
	if err != nil {
		http.Error(w, "invalid session", http.StatusUnauthorized)
		return
	}

	key := fmt.Sprintf("bookmarks:%s", username)
	if err := rdb2.RPush(ctx2, key, req.URL).Err(); err != nil {
		http.Error(w, "redis error", http.StatusInternalServerError)
		return
	}

	w.Write([]byte("bookmark added"))

}
