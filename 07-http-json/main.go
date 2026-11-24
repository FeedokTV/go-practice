package httpjson

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type UserCreateRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type User struct {
	ID           string `json:"id"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	passwordHash string `json:"-"`
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	r.Body = http.MaxBytesReader(w, r.Body, 1<<20) // limit: 1MB
	defer r.Body.Close()

	var input UserCreateRequest

	dec := json.NewDecoder(r.Body)
	if err := dec.Decode(&input); err != nil {
		var maxErr *http.MaxBytesError
		var synErr *json.SyntaxError
		var typeErr *json.UnmarshalTypeError

		switch {
		case errors.As(err, &maxErr):
			http.Error(w, "Request body too large", http.StatusRequestEntityTooLarge)
		case errors.As(err, &synErr), errors.Is(err, io.ErrUnexpectedEOF):
			http.Error(w, "Malformed JSON", http.StatusBadRequest)
		case errors.As(err, &typeErr):
			http.Error(w, "Invalid value for field "+typeErr.Field, http.StatusBadRequest)
		case errors.Is(err, io.EOF):
			http.Error(w, "Request body must not be empty", http.StatusBadRequest)
		default:
			http.Error(w, "Bad request", http.StatusBadRequest)
		}
		return // always stop on error

	}

	if err := dec.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, "Request body should contain only one JSON object", http.StatusBadRequest)
		return
	}

	if input.Username == "" || input.Email == "" || input.Password == "" {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	id, err := uuid.NewV7()
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	hashedPass, err := hashPassword(input.Password)
	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	user := User{
		ID:           id.String(),
		Username:     input.Username,
		Email:        input.Email,
		passwordHash: hashedPass,
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	if err := json.NewEncoder(w).Encode(user); err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
}
