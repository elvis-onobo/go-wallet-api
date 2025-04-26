package handler

import (
	"encoding/json"
	"net/http"
	"strings"

	"github.com/elvis-onobo/go-wallet-api/internal/http/middleware"
	"github.com/elvis-onobo/go-wallet-api/pkg/db"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID       int     `json:"id"`
	Username string  `json:"username"`
	Password string  `json:"password,omitempty"`
	Balance  float64 `json:"balance"`
}

var users = []User{}

func SignupHandler(w http.ResponseWriter, r *http.Request) {
	var user User

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)

	var id int
	err := db.Conn.QueryRow("INSERT INTO users (username, password) VALUES ($1, $2) RETURNING id", user.Username, string(hashedPassword)).Scan(&id)

	if err != nil {
		http.Error(w, "Could not create user", http.StatusInternalServerError)
		return
	}

	user.Password = ""
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(map[string]interface{}{
		"id":       id,
		"username": user.Username,
	})
}

func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	var foundUser *User
	for _, u := range users {
		if strings.EqualFold(u.Username, req.Username) {
			foundUser = &u
			break
		}
	}

	if foundUser == nil {
		http.Error(w, "User not found", http.StatusUnauthorized)
		return
	}

	err := bcrypt.CompareHashAndPassword([]byte(foundUser.Password), []byte(req.Password))
	if err != nil {
		http.Error(w, "Invalid password", http.StatusUnauthorized)
		return
	}

	token, err := generateJWT(foundUser.ID)
	if err != nil {
		http.Error(w, "Error generating token", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]string{
		"token":    token,
		"username": req.Username,
	})
}

func FundWalletHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(int)

	for i, user := range users {
		if user.ID == userID {
			users[i].Balance += req.Amount
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users[i])
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}

func WithdrawHandler(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Amount float64 `json:"amount"`
	}

	if err := json.NewDecoder(r.Body).Decode((&req)); err != nil {
		http.Error(w, "Invalid Request", http.StatusBadRequest)
		return
	}

	if req.Amount <= 0 {
		http.Error(w, "Amount must be greater than zero", http.StatusBadRequest)
		return
	}

	userID := r.Context().Value(middleware.UserIDKey).(int)

	for i, user := range users {
		if user.ID == userID {
			if user.Balance < req.Amount {
				http.Error(w, "Insufficient funds", http.StatusBadRequest)
				return
			}

			users[i].Balance -= req.Amount
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(users[i])
			return
		}
	}

	http.Error(w, "User not found", http.StatusNotFound)
}
