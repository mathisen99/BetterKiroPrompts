package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleKickoffGenerate_ValidInput(t *testing.T) {
	body := `{"answers":{"projectIdentity":"Test project","successCriteria":"Done when working"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/kickoff/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleKickoffGenerate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp KickoffResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if resp.Prompt == "" {
		t.Error("Expected non-empty prompt")
	}
}

func TestHandleKickoffGenerate_MissingProjectIdentity(t *testing.T) {
	body := `{"answers":{"successCriteria":"Done when working"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/kickoff/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleKickoffGenerate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestHandleKickoffGenerate_MalformedJSON(t *testing.T) {
	body := `{invalid json}`
	req := httptest.NewRequest(http.MethodPost, "/api/kickoff/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleKickoffGenerate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestHandleKickoffGenerate_WrongMethod(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/api/kickoff/generate", nil)
	w := httptest.NewRecorder()

	HandleKickoffGenerate(w, req)

	if w.Code != http.StatusMethodNotAllowed {
		t.Errorf("Expected 405, got %d", w.Code)
	}
}
