package models

import (
	"errors"
	"strings"
)

// Item represents a basic item in the system
type Item struct {
	ID          int    `json:"id" db:"id"`
	Name        string `json:"name" db:"name"`
	Description string `json:"description" db:"description"`
	Status      string `json:"status" db:"status"`
}

var allowedStatuses = map[string]bool{
	"active":    true,
	"archived":  true,
	"completed": true,
}

// Validate checks if an item has valid data
func (i *Item) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("item name cannot be empty")
	}
	if len(strings.TrimSpace(i.Name)) > 255 {
		return errors.New("item name cannot exceed 255 characters")
	}
	if len(strings.TrimSpace(i.Description)) > 1024 {
		return errors.New("item description cannot exceed 1024 characters")
	}
	if strings.TrimSpace(i.Status) == "" {
		i.Status = "active"
	}
	if !allowedStatuses[i.Status] {
		return errors.New("item status must be one of: active, archived, completed")
	}
	return nil
}

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// Validate checks if a create request is valid
func (r *CreateItemRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("item name cannot be empty")
	}
	if len(strings.TrimSpace(r.Name)) > 255 {
		return errors.New("item name cannot exceed 255 characters")
	}
	return nil
}

// UpdateItemRequest represents the request body for updating an item
type UpdateItemRequest struct {
	Name        string `json:"name"`
	Description string `json:"description,omitempty"`
	Status      string `json:"status,omitempty"`
}

// Validate checks if an update request is valid
func (r *UpdateItemRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("item name cannot be empty")
	}
	if len(strings.TrimSpace(r.Name)) > 255 {
		return errors.New("item name cannot exceed 255 characters")
	}
	if len(strings.TrimSpace(r.Description)) > 1024 {
		return errors.New("item description cannot exceed 1024 characters")
	}
	if strings.TrimSpace(r.Status) == "" {
		r.Status = "active"
	}
	if !allowedStatuses[r.Status] {
		return errors.New("item status must be one of: active, archived, completed")
	}
	return nil
}
