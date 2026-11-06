package handler

import (
	"context"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestContact(t *testing.T) {
	valid := `{"name":"  Alex  ","email":"alex@example.com","organization":"Example","message":"Hello"}`
	tests := []struct {
		name, method, body string
		sendErr            error
		want               int
		called             bool
	}{
		{"valid", http.MethodPost, valid, nil, http.StatusOK, true},
		{"delivery failure", http.MethodPost, valid, errors.New("unavailable"), http.StatusServiceUnavailable, true},
		{"wrong method", http.MethodGet, "", nil, http.StatusMethodNotAllowed, false},
		{"missing name", http.MethodPost, `{"email":"alex@example.com","message":"Hello"}`, nil, http.StatusBadRequest, false},
		{"invalid email", http.MethodPost, `{"name":"Alex","email":"not an email","message":"Hello"}`, nil, http.StatusBadRequest, false},
		{"missing message", http.MethodPost, `{"name":"Alex","email":"alex@example.com","message":" "}`, nil, http.StatusBadRequest, false},
		{"unknown field", http.MethodPost, `{"name":"Alex","email":"alex@example.com","message":"Hi","extra":true}`, nil, http.StatusBadRequest, false},
		{"multiple objects", http.MethodPost, valid + `{}`, nil, http.StatusBadRequest, false},
		{"oversized body", http.MethodPost, strings.Repeat("x", 17000), nil, http.StatusBadRequest, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			called := false
			recorder := httptest.NewRecorder()
			request := httptest.NewRequest(tt.method, "/api/contact", strings.NewReader(tt.body))
			handleContact(recorder, request, func(_ context.Context, message contactMessage) error {
				called = true
				if message.Name != "Alex" {
					t.Errorf("name = %q, want trimmed name", message.Name)
				}
				return tt.sendErr
			})
			if recorder.Code != tt.want || called != tt.called {
				t.Errorf("status = %d, called = %t; want %d, %t", recorder.Code, called, tt.want, tt.called)
			}
		})
	}
}

func TestSendContactRequiresConfiguration(t *testing.T) {
	t.Setenv("RESEND_API_KEY", "")
	if err := sendContact(context.Background(), contactMessage{}); err == nil {
		t.Fatal("expected an error when delivery is not configured")
	}
}
