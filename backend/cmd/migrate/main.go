package main

import (
	"database/sql"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sort"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	// Database connection string
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "postgres")
	dbPassword := getEnv("DB_PASSWORD", "root")
	dbName := getEnv("DB_NAME", "hotel_ecommerce")

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		dbHost, dbPort, dbUser, dbPassword, dbName)

	// Connect to database
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer db.Close()

	// Test connection
	if err := db.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}

	fmt.Println("Connected to PostgreSQL database successfully!")

	// Create migrations table if it doesn't exist
	createMigrationsTable(db)

	// Run migrations
	runMigrations(db)

	fmt.Println("Database migration completed successfully!")
}

func createMigrationsTable(db *sql.DB) {
	query := `
	CREATE TABLE IF NOT EXISTS migrations (
		id SERIAL PRIMARY KEY,
		filename VARCHAR(255) UNIQUE NOT NULL,
		executed_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	if err != nil {
		log.Fatal("Failed to create migrations table:", err)
	}
}

func runMigrations(db *sql.DB) {
	migrationsDir := "migrations"

	// Get list of migration files
	files, err := ioutil.ReadDir(migrationsDir)
	if err != nil {
		log.Fatal("Failed to read migrations directory:", err)
	}

	// Sort files by name to ensure correct order
	var migrationFiles []string
	for _, file := range files {
		if filepath.Ext(file.Name()) == ".sql" {
			migrationFiles = append(migrationFiles, file.Name())
		}
	}
	sort.Strings(migrationFiles)

	// Check which migrations have already been run
	executedMigrations := getExecutedMigrations(db)

	// Run pending migrations
	for _, filename := range migrationFiles {
		if _, exists := executedMigrations[filename]; !exists {
			fmt.Printf("Running migration: %s\n", filename)

			// Read migration file
			filePath := filepath.Join(migrationsDir, filename)
			content, err := ioutil.ReadFile(filePath)
			if err != nil {
				log.Fatalf("Failed to read migration file %s: %v", filename, err)
			}

			// Execute migration
			_, err = db.Exec(string(content))
			if err != nil {
				log.Fatalf("Failed to execute migration %s: %v", filename, err)
			}

			// Record migration as executed
			_, err = db.Exec("INSERT INTO migrations (filename) VALUES ($1)", filename)
			if err != nil {
				log.Fatalf("Failed to record migration %s: %v", filename, err)
			}

			fmt.Printf("Migration %s completed successfully!\n", filename)
		} else {
			fmt.Printf("Migration %s already executed, skipping...\n", filename)
		}
	}
}

func getExecutedMigrations(db *sql.DB) map[string]bool {
	executed := make(map[string]bool)

	rows, err := db.Query("SELECT filename FROM migrations")
	if err != nil {
		log.Fatal("Failed to get executed migrations:", err)
	}
	defer rows.Close()

	for rows.Next() {
		var filename string
		if err := rows.Scan(&filename); err != nil {
			log.Fatal("Failed to scan migration filename:", err)
		}
		executed[filename] = true
	}

	return executed
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
