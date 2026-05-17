package handlers

import (
	"bytes"
	"context"
	"encoding/json"
	"log/slog"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"CRUD_Project/database"
	"CRUD_Project/models"
)

func TestItemHandlerGetAll(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	// Create test items
	handler.repo.Create("Item 1")
	handler.repo.Create("Item 2")

	req, _ := http.NewRequest(http.MethodGet, "/api/items", nil)
	rr := httptest.NewRecorder()

	handler.GetAll(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response struct {
		Data []models.Item `json:"data"`
	}
	json.NewDecoder(rr.Body).Decode(&response)

	if len(response.Data) != 2 {
		t.Errorf("GetAll() returned %d items, want 2", len(response.Data))
	}
}

func TestItemHandlerGetByID(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	created, _ := handler.repo.Create("Test Item")

	req := httptest.NewRequest(http.MethodGet, "/api/items/1", nil)
	// Set the ID in the request context
	ctx := context.WithValue(req.Context(), idKey, "1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetByID(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response struct {
		Data models.Item `json:"data"`
	}
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.ID != created.ID {
		t.Errorf("GetByID() returned item ID %d, want %d", response.Data.ID, created.ID)
	}
}

func TestItemHandlerGetByIDNotFound(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	req := httptest.NewRequest(http.MethodGet, "/api/items/9999", nil)
	// Set the ID in the request context
	ctx := context.WithValue(req.Context(), idKey, "9999")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.GetByID(rr, req)

	if status := rr.Code; status != http.StatusNotFound {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNotFound)
	}
}

func TestItemHandlerCreate(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	reqBody := models.CreateItemRequest{Name: "New Item"}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if status := rr.Code; status != http.StatusCreated {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusCreated)
	}

	var response struct {
		Data models.Item `json:"data"`
	}
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.Name != "New Item" {
		t.Errorf("Create() returned item name %s, want %s", response.Data.Name, "New Item")
	}
}

func TestItemHandlerCreateInvalid(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	reqBody := models.CreateItemRequest{Name: ""}
	body, _ := json.Marshal(reqBody)

	req, _ := http.NewRequest(http.MethodPost, "/api/items", bytes.NewReader(body))
	rr := httptest.NewRecorder()

	handler.Create(rr, req)

	if status := rr.Code; status != http.StatusUnprocessableEntity {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusUnprocessableEntity)
	}
}

func TestItemHandlerUpdate(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	created, _ := handler.repo.Create("Original Name")

	reqBody := models.UpdateItemRequest{Name: "Updated Name"}
	body, _ := json.Marshal(reqBody)

	req := httptest.NewRequest(http.MethodPut, "/api/items/1", bytes.NewReader(body))
	// Set the ID in the request context
	ctx := context.WithValue(req.Context(), idKey, "1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.Update(rr, req)

	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	var response struct {
		Data models.Item `json:"data"`
	}
	json.NewDecoder(rr.Body).Decode(&response)

	if response.Data.ID != created.ID || response.Data.Name != "Updated Name" {
		t.Errorf("Update() failed")
	}
}

func TestItemHandlerDelete(t *testing.T) {
	handler, cleanup := setupTestHandler(t)
	defer cleanup()

	handler.repo.Create("Item to Delete")

	req := httptest.NewRequest(http.MethodDelete, "/api/items/1", nil)
	// Set the ID in the request context
	ctx := context.WithValue(req.Context(), idKey, "1")
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler.Delete(rr, req)

	if status := rr.Code; status != http.StatusNoContent {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusNoContent)
	}
}

// Helper functions
func setupTestHandler(t *testing.T) (*ItemHandler, func()) {
	repo, err := database.New("file::memory:?cache=shared&mode=rwc")
	if err != nil {
		t.Fatalf("failed to setup test repository: %v", err)
	}

	logger := slog.New(slog.NewTextHandler(os.Stdout, &slog.HandlerOptions{
		Level: slog.LevelError, // Suppress logs in tests
	}))

	handler := NewItemHandler(repo, logger)

	cleanup := func() {
		repo.Close()
	}

	return handler, cleanup
}
