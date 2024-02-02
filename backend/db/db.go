// db.go
package db

import (
	"database/sql"
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

// InitDB initializes the database connection
func InitDB() (*sql.DB, error) {
	// Load .env file
	if err := godotenv.Load(); err != nil {
		log.Fatal("Error loading .env file")
	}

	// Get environment variables
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

	// Build the connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open a DB connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err // Return an error if the connection fails
	}

	// Check the connection
	err = db.Ping()
	if err != nil {
		return nil, err // Return an error if the ping fails
	}

	execQuery := func(query string) error {
		_, err := db.Exec(query)
		if err != nil {
			return fmt.Errorf("error executing query: %v", err)
		}
		return nil
	}

	tables := []string{
		`CREATE TABLE IF NOT EXISTS ingredients (
            ingredient_id SERIAL PRIMARY KEY,
            ingredient_name VARCHAR(255) NOT NULL,
            ingredient_description TEXT
        );`,
		`CREATE TABLE IF NOT EXISTS recipes (
            recipe_id SERIAL PRIMARY KEY,
            recipe_name VARCHAR(255) NOT NULL,
            recipe_description TEXT,
            cook_time INT
        );`,
		`CREATE TABLE IF NOT EXISTS recipe_ingredients (
            recipe_ingredient_id SERIAL PRIMARY KEY,
            recipe_id INT NOT NULL,
            ingredient_id INT NOT NULL,
            quantity INT NOT NULL,
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id),
            FOREIGN KEY (ingredient_id) REFERENCES ingredients(ingredient_id)
        );`,
		`CREATE TABLE IF NOT EXISTS recipe_steps (
            recipe_step_id SERIAL PRIMARY KEY,
            recipe_id INT NOT NULL,
            step_number INT NOT NULL,
            step_description TEXT NOT NULL,
            FOREIGN KEY (recipe_id) REFERENCES recipes(recipe_id)
        );`,
	}

	for _, qry := range tables {
		if err := execQuery(qry); err != nil {
			return nil, err
		}
	}

	fmt.Println("Successfully connected!")
	return db, nil // Return the database instance and no error

	return db, nil

}
