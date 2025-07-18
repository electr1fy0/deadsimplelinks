package handler

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Bookmark struct {
	ID        int
	URL       string
	Title     string
	CreatedAt time.Time
}

type BookmarkStore struct {
	DB *sql.DB
}

func Handler(w http.ResponseWriter, r *http.Request) {
	db, err := connectDB()
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	store := &BookmarkStore{DB: db}

	bookmarks, err := store.GetAll()
	if err != nil {
		log.Fatal("fetch failed:", err)
	}

	fmt.Println("Bookmarks:")
	for _, b := range bookmarks {
		fmt.Fprintf(w, "%d | %s | %s | %s\n", b.ID, b.Title, b.URL, b.CreatedAt.Format("2006-01-02"))
	}

	// fmt.Fprint(w, "meow")

}

func (store *BookmarkStore) GetAll() ([]Bookmark, error) {
	rows, err := store.DB.Query("SELECT id, url, title, created_at FROM marks")
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bookmarks []Bookmark

	for rows.Next() {
		var b Bookmark
		if err := rows.Scan(&b.ID, &b.URL, &b.Title, &b.CreatedAt); err != nil {
			return nil, err
		}
		bookmarks = append(bookmarks, b)
	}
	return bookmarks, rows.Err()
}

func connectDB() (*sql.DB, error) {
	dsn := os.Getenv("SUPABASE_URL")
	db, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("connection failed: %w", err)
	}
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed: %w", err)
	}
	return db, nil
}
