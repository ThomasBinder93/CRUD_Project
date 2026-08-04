package database

import (
	"os"
	"testing"

	"CRUD_Project/models"
)

func TestRepositoryCreate(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	tests := []struct {
		name    string
		input   string
		wantErr bool
	}{
		{
			name:    "valid item",
			input:   "Test Item",
			wantErr: false,
		},
		{
			name:    "item with special chars",
			input:   "Item with !@#$%^&*()",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.Create(tt.input, "test description", "active")
			if (err != nil) != tt.wantErr {
				t.Errorf("Create() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && item.Name != tt.input {
				t.Errorf("Create() name = %v, want %v", item.Name, tt.input)
			}
		})
	}
}

func TestRepositoryGetAll(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	// Create test items
	repo.Create("Item 1", "First item", "active")
	repo.Create("Item 2", "Second item", "active")
	repo.Create("Item 3", "Third item", "active")

	items, err := repo.GetAll("")
	if err != nil {
		t.Errorf("GetAll() error = %v", err)
		return
	}

	if len(items) != 3 {
		t.Errorf("GetAll() returned %d items, want 3", len(items))
	}
}

func TestRepositoryGetByID(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	// Create a test item
	created, _ := repo.Create("Test Item", "Test item description", "active")

	tests := []struct {
		name    string
		id      int
		wantErr bool
		wantNil bool
	}{
		{
			name:    "existing item",
			id:      created.ID,
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "non-existing item",
			id:      9999,
			wantErr: false,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.GetByID(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (item == nil) != tt.wantNil {
				t.Errorf("GetByID() item nil = %v, wantNil %v", item == nil, tt.wantNil)
			}
		})
	}
}

func TestRepositoryUpdate(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	// Create a test item
	created, _ := repo.Create("Original Name", "Original description", "active")

	tests := []struct {
		name    string
		id      int
		newName string
		wantErr bool
		wantNil bool
	}{
		{
			name:    "update existing item",
			id:      created.ID,
			newName: "Updated Name",
			wantErr: false,
			wantNil: false,
		},
		{
			name:    "update non-existing item",
			id:      9999,
			newName: "Some Name",
			wantErr: false,
			wantNil: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			item, err := repo.Update(tt.id, tt.newName, "Updated description", "completed")
			if (err != nil) != tt.wantErr {
				t.Errorf("Update() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if (item == nil) != tt.wantNil {
				t.Errorf("Update() item nil = %v, wantNil %v", item == nil, tt.wantNil)
			}
			if !tt.wantNil && item.Name != tt.newName {
				t.Errorf("Update() name = %v, want %v", item.Name, tt.newName)
			}
		})
	}
}

func TestRepositoryDelete(t *testing.T) {
	repo := setupTestDB(t)
	defer cleanupTestDB(t, repo)

	// Create a test item
	created, _ := repo.Create("Item to Delete", "Deletion description", "active")

	tests := []struct {
		name    string
		id      int
		wantErr bool
	}{
		{
			name:    "delete existing item",
			id:      created.ID,
			wantErr: false,
		},
		{
			name:    "delete non-existing item",
			id:      9999,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.Delete(tt.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestModelValidation(t *testing.T) {
	tests := []struct {
		name    string
		item    *models.Item
		wantErr bool
	}{
		{
			name:    "valid item",
			item:    &models.Item{ID: 1, Name: "Valid Item"},
			wantErr: false,
		},
		{
			name:    "empty name",
			item:    &models.Item{ID: 1, Name: ""},
			wantErr: true,
		},
		{
			name:    "whitespace only",
			item:    &models.Item{ID: 1, Name: "   "},
			wantErr: true,
		},
		{
			name:    "name too long",
			item:    &models.Item{ID: 1, Name: string(make([]byte, 256))},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.item.Validate()
			if (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

// Helper functions
func setupTestDB(t *testing.T) *Repository {
	repo, err := New("file::memory:?cache=shared&mode=rwc")
	if err != nil {
		t.Fatalf("failed to setup test database: %v", err)
	}
	return repo
}

func cleanupTestDB(t *testing.T, repo *Repository) {
	if err := repo.Close(); err != nil {
		t.Errorf("failed to close test database: %v", err)
	}
	os.Remove("./test.db")
}
