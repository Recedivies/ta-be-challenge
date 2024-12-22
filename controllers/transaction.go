package controllers

import (
	"encoding/json"
	"errors"
	"net/http"
	"strings"

	"github.com/Recedivies/ta-be-challenge/db"
)

const (
	WITHDRAWAL string = "WITHDRAWAL"
	TOPUP      string = "TOPUP"
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

	err := ProcessTransaction(req.UserID, TOPUP, req.Amount)
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

	err := ProcessTransaction(req.UserID, WITHDRAWAL, req.Amount)
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

	adjustment := amount
	if transactionType == WITHDRAWAL {
		adjustment = -amount
	}

	_, err = tx.Exec(`
		UPDATE users
		SET balance = balance + $1
		WHERE id = $2
	`, adjustment, userID)
	if err != nil {
		if strings.Contains(err.Error(), "check constraint") {
			return errors.New("insufficient balance")
		}

		return errors.New("failed to update balance")
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
