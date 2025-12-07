package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
)

func main() {
	connStr := "host=65.21.6.209 port=5432 user=postgres password=Xk9mP2vL8nQ4wR7jT3yH6sF1zB5cN0aD dbname=mydb sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to database: %v\n", err)
		os.Exit(1)
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to ping database: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Connected to PostgreSQL successfully")

	schemaFile := "../database/schema.sql"
	schema, err := os.ReadFile(schemaFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read schema file: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Applying database schema...")

	if _, err := db.Exec(string(schema)); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to apply schema: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("Database schema applied successfully!")

	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM pg_tables WHERE schemaname = 'public'").Scan(&count); err != nil {
		fmt.Fprintf(os.Stderr, "Failed to verify tables: %v\n", err)
		os.Exit(1)
	}

	fmt.Printf("Verified: %d tables created in database\n", count)
}
