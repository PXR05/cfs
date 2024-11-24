package server

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"cfs/proc"
)

func setupTestServer(t *testing.T) *Server {
	s := &Server{}
	if err := s.Init(); err != nil {
		t.Fatalf("Failed to initialize server: %v", err)
	}
	return s
}

func TestEndpoints(t *testing.T) {
	s := setupTestServer(t)
	defer func() {
		if err := s.db.Cleanup(); err != nil {
			t.Errorf("Failed to cleanup test data: %v", err)
		}
		s.Close()
	}()

	// Create Categories
	t.Run("POST /cfs/c", func(t *testing.T) {
		categories := []proc.Category{
			{Name: "TestCategory1", Keywords: []string{"test1", "testing1"}},
			{Name: "TestCategory2", Keywords: []string{"test2", "testing2"}},
		}
		body, _ := json.Marshal(categories)
		req := httptest.NewRequest("POST", "/cfs/c", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		s.handleCreateCategories(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	// Get Categories
	t.Run("GET /cfs/c", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cfs/c", nil)
		w := httptest.NewRecorder()
		s.handleGetCategories(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response proc.CategoryOutputData
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if len(response.Categories) < 2 {
			t.Error("Expected at least 2 categories")
		}
	})

	// Create Classifications
	t.Run("POST /cfs/i", func(t *testing.T) {
		input := proc.InputData{
			Items: []string{"test item 1", "test item 2"},
		}
		body, _ := json.Marshal(input)
		req := httptest.NewRequest("POST", "/cfs/i", bytes.NewBuffer(body))
		w := httptest.NewRecorder()
		s.handleCreateClassifications(w, req)

		if w.Code != http.StatusCreated {
			t.Errorf("Expected status %d, got %d", http.StatusCreated, w.Code)
		}
	})

	// Get Classifications
	t.Run("GET /cfs/i", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cfs/i", nil)
		w := httptest.NewRecorder()
		s.handleGetClassifications(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}

		var response proc.ClassificationOutputData
		if err := json.NewDecoder(w.Body).Decode(&response); err != nil {
			t.Fatalf("Failed to decode response: %v", err)
		}
		if len(response.Results) < 2 {
			t.Error("Expected at least 2 classifications")
		}
	})

	// Get Single Category
	t.Run("GET /cfs/c/{category}", func(t *testing.T) {
		req := httptest.NewRequest("GET", "/cfs/c?category=TestCategory1", nil)
		w := httptest.NewRecorder()
		s.handleGetCategory(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	// Get Single Classification
	t.Run("GET /cfs/i/{item}", func(t *testing.T) {
		itemQuery := url.QueryEscape("test item 1")
		req := httptest.NewRequest("GET", "/cfs/i?item="+itemQuery, nil)
		w := httptest.NewRecorder()
		s.handleGetClassification(w, req)

		if w.Code != http.StatusOK {
			t.Errorf("Expected status %d, got %d", http.StatusOK, w.Code)
		}
	})

	// Error Cases
	t.Run("Error Cases", func(t *testing.T) {
		// Missing category
		req := httptest.NewRequest("GET", "/cfs/c", nil)
		w := httptest.NewRecorder()
		s.handleGetCategory(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for missing category, got %d", http.StatusBadRequest, w.Code)
		}

		// Missing item
		req = httptest.NewRequest("GET", "/cfs/i", nil)
		w = httptest.NewRecorder()
		s.handleGetClassification(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for missing item, got %d", http.StatusBadRequest, w.Code)
		}

		// Invalid JSON in create classification
		req = httptest.NewRequest("POST", "/cfs/i", bytes.NewBufferString("invalid json"))
		w = httptest.NewRecorder()
		s.handleCreateClassifications(w, req)
		if w.Code != http.StatusBadRequest {
			t.Errorf("Expected status %d for invalid JSON, got %d", http.StatusBadRequest, w.Code)
		}
	})
}
