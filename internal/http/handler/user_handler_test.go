package handler_test

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	handler "github.com/elvis-onobo/go-wallet-api/internal/http/handler"
	"github.com/elvis-onobo/go-wallet-api/pkg/db"
)

func setupTestDB() {
	var err error
	db.Conn, err = sql.Open("postgres", "postgres://postgres:password@localhost:5434/walletapi?sslmode=disable")
	if err != nil {
		panic(err)
	}
}

func TestSignupHandler(t *testing.T) {
	setupTestDB()

	payload := map[string]string{
		"username": "testuser",
		"password": "secret",
	}
	body, _ := json.Marshal(payload)

	req := httptest.NewRequest(http.MethodPost, "/signup", bytes.NewReader(body))
	w := httptest.NewRecorder()

	handler.SignupHandler(w, req)

	// resp := w.Result()
	// if resp.StatusCode != http.StatusCreated {
	// 	t.Errorf("expected status 201 Created, got %d", resp.StatusCode)
	// }
}
