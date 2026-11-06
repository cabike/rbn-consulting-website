package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/mail"
	"os"
	"strings"
	"time"
)

type contactMessage struct {
	Name         string `json:"name"`
	Email        string `json:"email"`
	Organization string `json:"organization"`
	Message      string `json:"message"`
}

// Handler accepts contact requests for a Go HTTP serverless runtime.
func Handler(w http.ResponseWriter, r *http.Request) {
	handleContact(w, r, sendContact)
}

func handleContact(w http.ResponseWriter, r *http.Request, send func(context.Context, contactMessage) error) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Cache-Control", "no-store")
	if r.Method != http.MethodPost {
		w.Header().Set("Allow", http.MethodPost)
		http.Error(w, `{"error":"Method not allowed"}`, http.StatusMethodNotAllowed)
		return
	}

	r.Body = http.MaxBytesReader(w, r.Body, 16<<10)
	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()
	var message contactMessage
	if err := decoder.Decode(&message); err != nil {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}
	if err := decoder.Decode(&struct{}{}); err != io.EOF {
		http.Error(w, `{"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}
	message.Name = strings.TrimSpace(message.Name)
	message.Email = strings.TrimSpace(message.Email)
	message.Organization = strings.TrimSpace(message.Organization)
	message.Message = strings.TrimSpace(message.Message)
	address, err := mail.ParseAddress(message.Email)
	if message.Name == "" || len(message.Name) > 100 || len(message.Email) > 254 ||
		err != nil || address.Address != message.Email || len(message.Organization) > 150 ||
		message.Message == "" || len(message.Message) > 4000 {
		http.Error(w, `{"error":"Please check your details"}`, http.StatusBadRequest)
		return
	}
	if err := send(r.Context(), message); err != nil {
		http.Error(w, `{"error":"Unable to send message"}`, http.StatusServiceUnavailable)
		return
	}
	w.WriteHeader(http.StatusOK)
	_, _ = io.WriteString(w, `{"ok":true}`)
}

func sendContact(ctx context.Context, message contactMessage) error {
	key := os.Getenv("RESEND_API_KEY")
	from := os.Getenv("CONTACT_FROM_EMAIL")
	to := os.Getenv("CONTACT_TO_EMAIL")
	if key == "" || from == "" || to == "" {
		return errors.New("contact delivery is not configured")
	}
	payload, err := json.Marshal(map[string]any{
		"from":     from,
		"to":       []string{to},
		"reply_to": message.Email,
		"subject":  "New RBN Consulting website inquiry",
		"text": fmt.Sprintf("Name: %s\nEmail: %s\nOrganization: %s\n\n%s",
			message.Name, message.Email, message.Organization, message.Message),
	})
	if err != nil {
		return err
	}
	ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://api.resend.com/emails", bytes.NewReader(payload))
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+key)
	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return errors.New("contact delivery failed")
	}
	return nil
}
