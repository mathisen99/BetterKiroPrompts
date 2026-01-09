package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleHooksGenerate_ValidInput(t *testing.T) {
	body := `{"preset":"light","techStack":{"hasGo":true}}`
	req := httptest.NewRequest(http.MethodPost, "/api/hooks/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleHooksGenerate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp HooksResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp.Files) == 0 {
		t.Error("Expected files in response")
	}
}

func TestHandleHooksGenerate_MissingPreset(t *testing.T) {
	body := `{"techStack":{"hasGo":true}}`
	req := httptest.NewRequest(http.MethodPost, "/api/hooks/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleHooksGenerate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestHandleHooksGenerate_PresetVariations(t *testing.T) {
	lightBody := `{"preset":"light","techStack":{"hasGo":true}}`
	lightReq := httptest.NewRequest(http.MethodPost, "/api/hooks/generate", bytes.NewBufferString(lightBody))
	lightW := httptest.NewRecorder()
	HandleHooksGenerate(lightW, lightReq)

	var lightResp HooksResponse
	json.NewDecoder(lightW.Body).Decode(&lightResp)

	strictBody := `{"preset":"strict","techStack":{"hasGo":true}}`
	strictReq := httptest.NewRequest(http.MethodPost, "/api/hooks/generate", bytes.NewBufferString(strictBody))
	strictW := httptest.NewRecorder()
	HandleHooksGenerate(strictW, strictReq)

	var strictResp HooksResponse
	json.NewDecoder(strictW.Body).Decode(&strictResp)

	if len(strictResp.Files) <= len(lightResp.Files) {
		t.Errorf("Expected strict to have more files than light, got %d vs %d", len(strictResp.Files), len(lightResp.Files))
	}
}

func TestHandleHooksGenerate_MalformedJSON(t *testing.T) {
	body := `{invalid}`
	req := httptest.NewRequest(http.MethodPost, "/api/hooks/generate", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	HandleHooksGenerate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
