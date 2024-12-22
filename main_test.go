package main

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"

	"github.com/Recedivies/ta-be-challenge/controllers"
	"github.com/Recedivies/ta-be-challenge/db"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const LOCALHOST string = "localhost"

func TestTopUpHandler(t *testing.T) {
	db.Connect(LOCALHOST)
	db.DropDB()
	db.InitDB()
	db.SeedDB()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(userID int, transactionType string, amount float64) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Invalid request body",
			requestBody: "invalid-json",
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body\n",
		},
		{
			name: "Invalid top-up amount",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"amount":  -10,
			},
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Amount must be greater than zero\n",
		},
		{
			name: "Successful top-up",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"amount":  100,
			},
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}` + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			var reqBody []byte
			if body, ok := tc.requestBody.(string); ok && body == "invalid-json" {
				reqBody = []byte(body)
			} else {
				reqBody, _ = json.Marshal(tc.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/topup", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			controllers.TopUpHandler(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if rr.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

func TestWithdrawalHandler(t *testing.T) {
	db.Connect(LOCALHOST)
	db.DropDB()
	db.InitDB()
	db.SeedDB()

	tests := []struct {
		name           string
		requestBody    interface{}
		mockFunc       func(userID int, transactionType string, amount float64) error
		expectedStatus int
		expectedBody   string
	}{
		{
			name:        "Invalid request body",
			requestBody: "invalid-json",
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Invalid request body\n",
		},
		{
			name: "Invalid withdrawal amount",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"amount":  -10,
			},
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Amount must be greater than zero\n",
		},
		{
			name: "Successful withdrawal",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"amount":  100,
			},
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusOK,
			expectedBody:   `{"status":"success"}` + "\n",
		},
		{
			name: "Insufficient balance",
			requestBody: map[string]interface{}{
				"user_id": 1,
				"amount":  9999999999,
			},
			mockFunc: func(userID int, transactionType string, amount float64) error {
				return nil
			},
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   `insufficient balance` + "\n",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {

			var reqBody []byte
			if body, ok := tc.requestBody.(string); ok && body == "invalid-json" {
				reqBody = []byte(body)
			} else {
				reqBody, _ = json.Marshal(tc.requestBody)
			}

			req := httptest.NewRequest(http.MethodPost, "/api/v1/withdrawal", bytes.NewReader(reqBody))
			req.Header.Set("Content-Type", "application/json")

			rr := httptest.NewRecorder()

			controllers.WithdrawalHandler(rr, req)

			if rr.Code != tc.expectedStatus {
				t.Errorf("expected status %d, got %d", tc.expectedStatus, rr.Code)
			}

			if rr.Body.String() != tc.expectedBody {
				t.Errorf("expected body %q, got %q", tc.expectedBody, rr.Body.String())
			}
		})
	}
}

const DEFAULT_USER_ID = 1

func TestConcurrency(t *testing.T) {
	db.Connect(LOCALHOST)
	db.DropDB()
	db.InitDB()
	db.SeedDB()

	var wg sync.WaitGroup
	const numRequests = 100
	errors := make(chan error, numRequests)

	for i := 0; i < numRequests; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			checkoutAmount := 100.00
			req := controllers.UserRequest{
				UserID: DEFAULT_USER_ID,
				Amount: checkoutAmount,
			}

			err := controllers.ProcessTransaction(req.UserID, "TOPUP", req.Amount)
			if err != nil {
				errors <- err
			}
		}()
	}

	wg.Wait()
	close(errors)

	for err := range errors {
		if err != nil {
			t.Fatalf("Error occurred during concurrent request: %v", err)
		}
	}

	var finalBalance string
	err := db.DB.QueryRow("SELECT balance FROM users WHERE id = $1", DEFAULT_USER_ID).Scan(&finalBalance)
	require.NoError(t, err)

	// Ahmadhi's account with ID=1 has 1000 balance, the result should be 1000+100 amount * 100 num of request = 11000
	assert.Equal(t, "11000.00", finalBalance)
}
