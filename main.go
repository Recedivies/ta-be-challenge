package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/Recedivies/ta-be-challenge/controllers"
	"github.com/Recedivies/ta-be-challenge/db"
)

func main() {
	db.Connect()

	log.Println("Dropping and reinitializing the database...")
	db.DropDB()
	db.InitDB()
	db.SeedDB()
	log.Println("Database reset successfully!")

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		response := map[string]string{
			"message": "Pong",
		}
		if err := json.NewEncoder(w).Encode(response); err != nil {
			http.Error(w, "Unable to encode response", http.StatusInternalServerError)
		}
	})

	http.HandleFunc("/api/v1/topup", controllers.TopUpHandler)
	http.HandleFunc("/api/v1/withdrawal", controllers.WithdrawalHandler)

	log.Println("Server running on http://localhost:9090")
	log.Fatal(http.ListenAndServe(":9090", nil))
}
