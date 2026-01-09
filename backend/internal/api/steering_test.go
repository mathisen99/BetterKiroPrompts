package api

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHandleSteeringGenerate_ValidInput(t *testing.T) {
	body := `{"config":{"projectName":"Test Project","projectDescription":"A test"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/steering/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleSteeringGenerate(w, req)

	if w.Code != http.StatusOK {
		t.Errorf("Expected 200, got %d: %s", w.Code, w.Body.String())
	}

	var resp SteeringResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("Failed to decode response: %v", err)
	}
	if len(resp.Files) == 0 {
		t.Error("Expected files in response")
	}
}

func TestHandleSteeringGenerate_MissingProjectName(t *testing.T) {
	body := `{"config":{"projectDescription":"A test"}}`
	req := httptest.NewRequest(http.MethodPost, "/api/steering/generate", bytes.NewBufferString(body))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	HandleSteeringGenerate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}

func TestHandleSteeringGenerate_ConditionalToggle(t *testing.T) {
	bodyWithout := `{"config":{"projectName":"Test"}}`
	reqWithout := httptest.NewRequest(http.MethodPost, "/api/steering/generate", bytes.NewBufferString(bodyWithout))
	wWithout := httptest.NewRecorder()
	HandleSteeringGenerate(wWithout, reqWithout)

	var respWithout SteeringResponse
	json.NewDecoder(wWithout.Body).Decode(&respWithout)

	bodyWith := `{"config":{"projectName":"Test","includeConditional":true}}`
	reqWith := httptest.NewRequest(http.MethodPost, "/api/steering/generate", bytes.NewBufferString(bodyWith))
	wWith := httptest.NewRecorder()
	HandleSteeringGenerate(wWith, reqWith)

	var respWith SteeringResponse
	json.NewDecoder(wWith.Body).Decode(&respWith)

	if len(respWith.Files) <= len(respWithout.Files) {
		t.Errorf("Expected more files with conditional, got %d vs %d", len(respWith.Files), len(respWithout.Files))
	}
}

func TestHandleSteeringGenerate_MalformedJSON(t *testing.T) {
	body := `{invalid}`
	req := httptest.NewRequest(http.MethodPost, "/api/steering/generate", bytes.NewBufferString(body))
	w := httptest.NewRecorder()

	HandleSteeringGenerate(w, req)

	if w.Code != http.StatusBadRequest {
		t.Errorf("Expected 400, got %d", w.Code)
	}
}
