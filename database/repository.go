package database

import (
	"database/sql"
	"errors"
	"fmt"
	"strings"
	"sync"

	"CRUD_Project/models"

	_ "modernc.org/sqlite"
)

// Repository handles all database operations
type Repository struct {
	db *sql.DB
	mu sync.RWMutex
}

// New creates a new repository and initializes the database
func New(dsn string) (*Repository, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("failed to open database: %w", err)
	}

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("failed to connect to database: %w", err)
	}

	repo := &Repository{db: db}

	// Initialize schema
	if err := repo.initSchema(); err != nil {
		return nil, fmt.Errorf("failed to initialize schema: %w", err)
	}

	return repo, nil
}

// initSchema creates tables if they don't exist
func (r *Repository) initSchema() error {
	createTableSQL := `CREATE TABLE IF NOT EXISTS items (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		name TEXT NOT NULL,
		description TEXT,
		status TEXT NOT NULL DEFAULT 'active'
	);`

	_, err := r.db.Exec(createTableSQL)
	if err != nil {
		return fmt.Errorf("failed to create items table: %w", err)
	}

	cols, err := r.getTableColumns("items")
	if err != nil {
		return err
	}

	if !cols["description"] {
		if _, err := r.db.Exec("ALTER TABLE items ADD COLUMN description TEXT"); err != nil {
			return fmt.Errorf("failed to add description column: %w", err)
		}
	}

	if !cols["status"] {
		if _, err := r.db.Exec("ALTER TABLE items ADD COLUMN status TEXT NOT NULL DEFAULT 'active'"); err != nil {
			return fmt.Errorf("failed to add status column: %w", err)
		}
	}

	return nil
}

func (r *Repository) getTableColumns(tableName string) (map[string]bool, error) {
	rows, err := r.db.Query(fmt.Sprintf("PRAGMA table_info(%s)", tableName))
	if err != nil {
		return nil, fmt.Errorf("failed to query table info: %w", err)
	}
	defer rows.Close()

	cols := make(map[string]bool)
	var cid int
	var name, ctype string
	var notnull, pk int
	var dfltValue sql.NullString

	for rows.Next() {
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return nil, fmt.Errorf("failed to scan table info: %w", err)
		}
		cols[name] = true
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error reading table info: %w", err)
	}

	return cols, nil
}

// GetAll retrieves all items from the database, optionally filtered by a search term.
func (r *Repository) GetAll(search string) ([]models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	query := "SELECT id, name, description, status FROM items"
	args := []interface{}{}
	if search != "" {
		query += " WHERE LOWER(name) LIKE ? OR LOWER(description) LIKE ?"
		pattern := "%" + strings.ToLower(search) + "%"
		args = append(args, pattern, pattern)
	}
	query += " ORDER BY id"

	rows, err := r.db.Query(query, args...)
	if err != nil {
		return nil, fmt.Errorf("failed to query items: %w", err)
	}
	defer rows.Close()

	var items []models.Item
	for rows.Next() {
		var item models.Item
		var description sql.NullString
		if err := rows.Scan(&item.ID, &item.Name, &description, &item.Status); err != nil {
			return nil, fmt.Errorf("failed to scan item: %w", err)
		}
		if description.Valid {
			item.Description = description.String
		}
		items = append(items, item)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during iteration: %w", err)
	}

	if items == nil {
		items = []models.Item{}
	}

	return items, nil
}

// GetByID retrieves a single item by ID
func (r *Repository) GetByID(id int) (*models.Item, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var item models.Item
	var description sql.NullString
	err := r.db.QueryRow("SELECT id, name, description, status FROM items WHERE id = ?", id).
		Scan(&item.ID, &item.Name, &description, &item.Status)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Item not found
		}
		return nil, fmt.Errorf("failed to get item: %w", err)
	}

	if description.Valid {
		item.Description = description.String
	}

	return &item, nil
}

// Create inserts a new item and returns it with ID
func (r *Repository) Create(name, description, status string) (*models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	if status == "" {
		status = "active"
	}

	result, err := r.db.Exec("INSERT INTO items (name, description, status) VALUES (?, ?, ?)", name, description, status)
	if err != nil {
		return nil, fmt.Errorf("failed to create item: %w", err)
	}

	id, err := result.LastInsertId()
	if err != nil {
		return nil, fmt.Errorf("failed to get last insert id: %w", err)
	}

	return &models.Item{
		ID:          int(id),
		Name:        name,
		Description: description,
		Status:      status,
	}, nil
}

// Update modifies an existing item
func (r *Repository) Update(id int, name, description, status string) (*models.Item, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	// Check if item exists
	var exists bool
	err := r.db.QueryRow("SELECT 1 FROM items WHERE id = ?", id).Scan(&exists)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil // Item not found
		}
		return nil, fmt.Errorf("failed to check item existence: %w", err)
	}

	result, err := r.db.Exec("UPDATE items SET name = ?, description = ?, status = ? WHERE id = ?", name, description, status, id)
	if err != nil {
		return nil, fmt.Errorf("failed to update item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return nil, fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil, nil // Item not found
	}

	return &models.Item{
		ID:          id,
		Name:        name,
		Description: description,
		Status:      status,
	}, nil
}

// Delete removes an item by ID
func (r *Repository) Delete(id int) error {
	r.mu.Lock()
	defer r.mu.Unlock()

	result, err := r.db.Exec("DELETE FROM items WHERE id = ?", id)
	if err != nil {
		return fmt.Errorf("failed to delete item: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected: %w", err)
	}

	if rowsAffected == 0 {
		return nil // Item not found, but no error
	}

	return nil
}

// Close closes the database connection
func (r *Repository) Close() error {
	return r.db.Close()
}
