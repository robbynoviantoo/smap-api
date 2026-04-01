package config

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

var DB *sql.DB

func InitDB() {
	dsn := App.DSN()
	db, err := sql.Open("mysql", dsn)
	if err != nil {
		log.Fatalf("[DB] Failed to open connection: %v", err)
	}

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatalf("[DB] Failed to ping: %v", err)
	}

	// Connection pool settings
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(10)

	DB = db
	log.Println("[DB] MySQL connected successfully")
}
