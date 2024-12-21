package controllers

import (
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/Recedivies/ta-be-challenge/db"
)

type UserRequest struct {
	UserID int64   `json:"user_id"`
	Amount float64 `json:"amount"`
}

func TopUpHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	err := ProcessTransaction(req.UserID, "TOPUP", req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func WithdrawalHandler(w http.ResponseWriter, r *http.Request) {
	var req UserRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	err := ProcessTransaction(req.UserID, "WITHDRAWAL", req.Amount)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{
		"status": "success",
	})
}

func ProcessTransaction(userID int64, transactionType string, amount float64) error {
	tx, err := db.DB.Begin()
	if err != nil {
		return errors.New("failed to start transaction")
	}
	defer tx.Rollback()

	var currentBalance float64
	err = tx.QueryRow("SELECT balance FROM users WHERE id = $1 FOR UPDATE", userID).Scan(&currentBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return errors.New("user not found")
		}
		return errors.New("failed to fetch user balance")
	}

	if transactionType == "WITHDRAWAL" {
		if currentBalance < amount {
			return errors.New("insufficient balance")
		}
		currentBalance -= amount
	} else if transactionType == "TOPUP" {
		currentBalance += amount
	}

	_, err = tx.Exec("UPDATE users SET balance = $1 WHERE id = $2", currentBalance, userID)
	if err != nil {
		return errors.New("failed to update user balance")
	}

	_, err = tx.Exec(`
		INSERT INTO transactions (user_id, transaction_type, amount, status)
		VALUES ($1, $2, $3, 'COMPLETED')
	`, userID, transactionType, amount)
	if err != nil {
		return errors.New("failed to record transaction")
	}

	if err := tx.Commit(); err != nil {
		return errors.New("failed to commit transaction")
	}

	return nil
}
