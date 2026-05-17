package models

import (
	"errors"
	"strings"
)

// Item represents a basic item in the system
type Item struct {
	ID   int    `json:"id" db:"id"`
	Name string `json:"name" db:"name"`
}

// Validate checks if an item has valid data
func (i *Item) Validate() error {
	if strings.TrimSpace(i.Name) == "" {
		return errors.New("item name cannot be empty")
	}
	if len(strings.TrimSpace(i.Name)) > 255 {
		return errors.New("item name cannot exceed 255 characters")
	}
	return nil
}

// CreateItemRequest represents the request body for creating an item
type CreateItemRequest struct {
	Name string `json:"name"`
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
	Name string `json:"name"`
}

// Validate checks if an update request is valid
func (r *UpdateItemRequest) Validate() error {
	if strings.TrimSpace(r.Name) == "" {
		return errors.New("item name cannot be empty")
	}
	if len(strings.TrimSpace(r.Name)) > 255 {
		return errors.New("item name cannot exceed 255 characters")
	}
	return nil
}
