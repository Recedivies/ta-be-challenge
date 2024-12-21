package db

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func Connect() {
	var err error
	connStr := "postgresql://root:secret@postgres:5432/finance?sslmode=disable"
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Fatalf("Failed to connect to the database: %v", err)
	}

	if err := DB.Ping(); err != nil {
		log.Fatalf("Failed to ping the database: %v", err)
	}
	log.Println("Connected to the database!")
}

func InitDB() {
	createUsersTable := `
		CREATE TABLE IF NOT EXISTS users (
			id SERIAL PRIMARY KEY,
			name VARCHAR(100) NOT NULL,
			balance NUMERIC(20, 2) NOT NULL DEFAULT 0
		);
	`

	createTransactionsTable := `
		CREATE TABLE IF NOT EXISTS transactions (
			id SERIAL PRIMARY KEY,
			user_id INT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
			transaction_type VARCHAR(20) NOT NULL,
			amount NUMERIC(20, 2) NOT NULL,
			status VARCHAR(20) NOT NULL,
			created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
		);
	`

	log.Println("Initializing database...")

	_, err := DB.Exec(createUsersTable)
	if err != nil {
		log.Fatalf("Failed to create users table: %v", err)
	}

	_, err = DB.Exec(createTransactionsTable)
	if err != nil {
		log.Fatalf("Failed to create transactions table: %v", err)
	}

	log.Println("Database tables created successfully!")
}

func DropDB() {
	dropTransactionsTable := "DROP TABLE IF EXISTS transactions;"
	dropUsersTable := "DROP TABLE IF EXISTS users;"

	log.Println("Dropping database tables...")

	_, err := DB.Exec(dropTransactionsTable)
	if err != nil {
		log.Fatalf("Failed to drop transactions table: %v", err)
	}

	_, err = DB.Exec(dropUsersTable)
	if err != nil {
		log.Fatalf("Failed to drop users table: %v", err)
	}

	log.Println("Database tables dropped successfully!")
}

func SeedDB() {
	log.Println("Seeding database...")

	// Seed Users
	insertUsers := `
		INSERT INTO users (name, balance) VALUES
		('Ahmadhi', 1000.00),
		('Prananta', 500.00),
		('Hastiputra', 200.00);
	`

	rows, err := DB.Query(insertUsers)
	if err != nil {
		log.Fatalf("Failed to seed users: %v", err)
	}
	defer rows.Close()

	log.Println("Database seeding completed successfully!")
}
