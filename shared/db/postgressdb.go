package db

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/lib/pq"
)

// Open opens & verifies a Postgres connection.
func Open(connStr string) *sql.DB {
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("❌ open db: %v", err)
	}
	if err = db.Ping(); err != nil {
		log.Fatalf("❌ ping db: %v", err)
	}

	// ← Add this line so you see it on startup:
	fmt.Println("✅ Database connection verified.")

	return db
}
