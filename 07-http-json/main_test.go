package httpjson

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestCreateUser(t *testing.T) {
	input := UserCreateRequest{
		Username: "test",
		Email:    "test@test.com",
		Password: "1234qwerty",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("unexpected error on marshalling: %v", err)
	}

	req := httptest.NewRequest("POST", "http://example/api/user/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		t.Fatalf("unexpected status, want: 201, got: %d", resp.StatusCode)
	}

	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if user.ID == "" {
		t.Errorf("expected not nil id")
	}

	if user.Username != "test" {
		t.Errorf("want username: test, got %s", user.Username)
	}

	if user.Email != "test@test.com" {
		t.Errorf("want email: test@test.com, got %s", user.Email)
	}

	if strings.Contains(strings.ToLower(w.Body.String()), "password") {
		t.Errorf("expected no password in response")
	}
}

func TestCreateUserInvalidInput(t *testing.T) {
	input := UserCreateRequest{
		Username: "test",
	}

	body, err := json.Marshal(input)
	if err != nil {
		t.Fatalf("unexpected error on marshalling: %v", err)
	}

	req := httptest.NewRequest("POST", "http://example/api/user/create", bytes.NewReader(body))
	w := httptest.NewRecorder()

	CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status, want: 400, got: %d", resp.StatusCode)
	}

	bodyString := strings.TrimSpace(w.Body.String())

	if bodyString != "Invalid request" {
		t.Errorf("expected error message in body, got %s", bodyString)
	}
}

func TestCreateUserInvalidRequestTwoJSON(t *testing.T) {

	body := `{"username":"test", "email": "test@test.com", "password":"123"},{}`

	req := httptest.NewRequest("POST", "http://example/api/user/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status, want: 400, got: %d", resp.StatusCode)
	}

	bodyString := strings.TrimSpace(w.Body.String())

	if bodyString != "Request body should contain only one JSON object" {
		t.Errorf("expected error message in body, got %s", bodyString)
	}
}

func TestCreateUserInvalidRequestBrokenObject(t *testing.T) {

	body := `{"username":`

	req := httptest.NewRequest("POST", "http://example/api/user/create", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	CreateUser(w, req)

	resp := w.Result()
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusBadRequest {
		t.Fatalf("unexpected status, want: 400, got: %d", resp.StatusCode)
	}

	bodyString := strings.TrimSpace(w.Body.String())

	if bodyString != "Malformed JSON" {
		t.Errorf("expected error message in body, got %s", bodyString)
	}
}

func TestLoopVariable(t *testing.T) {

	var prevPointer *int
	for i := range 5 {
		if i == 0 {
			prevPointer = &i
			continue
		}

		if prevPointer == &i {
			t.Errorf("%d is have the same pointer as previous", i)
		}
	}
}
