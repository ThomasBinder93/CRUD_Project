package handlers

import (
	"context"
	"errors"
	"io"
	"log/slog"
	"net/http"
	"strconv"

	"CRUD_Project/database"
	"CRUD_Project/models"
	"CRUD_Project/utils"
)

// itemIDKey is the context key for storing the item ID
type itemIDKey struct{}

var idKey = itemIDKey{}

// ItemHandler handles all item-related HTTP requests
type ItemHandler struct {
	repo   *database.Repository
	logger *slog.Logger
}

// NewItemHandler creates a new item handler
func NewItemHandler(repo *database.Repository, logger *slog.Logger) *ItemHandler {
	return &ItemHandler{
		repo:   repo,
		logger: logger,
	}
}

// GetAll retrieves all items
// @Summary Get all items
// @Description Retrieve a list of all items
// @Tags items
// @Accept json
// @Produce json
// @Success 200 {object} utils.SuccessResponse
// @Router /api/items [get]
func (h *ItemHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	items, err := h.repo.GetAll()
	if err != nil {
		h.logger.Error("failed to get items", slog.String("error", err.Error()))
		_ = utils.WriteError(w, utils.ErrorInternalServer("failed to retrieve items"))
		return
	}

	_ = utils.WriteSuccess(w, http.StatusOK, items)
}

// GetByID retrieves a single item by ID
// @Summary Get item by ID
// @Description Retrieve a specific item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 200 {object} utils.SuccessResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /api/items/{id} [get]
func (h *ItemHandler) GetByID(w http.ResponseWriter, r *http.Request) {
	idStr := r.Context().Value(idKey).(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = utils.WriteError(w, utils.ErrorBadRequest("invalid item ID"))
		return
	}

	item, err := h.repo.GetByID(id)
	if err != nil {
		h.logger.Error("failed to get item", slog.String("error", err.Error()))
		_ = utils.WriteError(w, utils.ErrorInternalServer("failed to retrieve item"))
		return
	}

	if item == nil {
		_ = utils.WriteError(w, utils.ErrorNotFound("item not found"))
		return
	}

	_ = utils.WriteSuccess(w, http.StatusOK, item)
}

// Create creates a new item
// @Summary Create a new item
// @Description Create a new item with the provided name
// @Tags items
// @Accept json
// @Produce json
// @Param request body models.CreateItemRequest true "Item data"
// @Success 201 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /api/items [post]
func (h *ItemHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req models.CreateItemRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		_ = utils.WriteError(w, utils.ErrorBadRequest("invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		_ = utils.WriteError(w, utils.ErrorUnprocessableEntity(err.Error()))
		return
	}

	item, err := h.repo.Create(req.Name)
	if err != nil {
		h.logger.Error("failed to create item", slog.String("error", err.Error()))
		_ = utils.WriteError(w, utils.ErrorInternalServer("failed to create item"))
		return
	}

	_ = utils.WriteSuccess(w, http.StatusCreated, item)
}

// Update updates an existing item
// @Summary Update an item
// @Description Update an existing item with new data
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Param request body models.UpdateItemRequest true "Item data"
// @Success 200 {object} utils.SuccessResponse
// @Failure 400 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 422 {object} utils.ErrorResponse
// @Router /api/items/{id} [put]
func (h *ItemHandler) Update(w http.ResponseWriter, r *http.Request) {
	idStr := r.Context().Value(idKey).(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = utils.WriteError(w, utils.ErrorBadRequest("invalid item ID"))
		return
	}

	var req models.UpdateItemRequest
	if err := decodeJSON(r.Body, &req); err != nil {
		_ = utils.WriteError(w, utils.ErrorBadRequest("invalid request body"))
		return
	}

	if err := req.Validate(); err != nil {
		_ = utils.WriteError(w, utils.ErrorUnprocessableEntity(err.Error()))
		return
	}

	item, err := h.repo.Update(id, req.Name)
	if err != nil {
		h.logger.Error("failed to update item", slog.String("error", err.Error()))
		_ = utils.WriteError(w, utils.ErrorInternalServer("failed to update item"))
		return
	}

	if item == nil {
		_ = utils.WriteError(w, utils.ErrorNotFound("item not found"))
		return
	}

	_ = utils.WriteSuccess(w, http.StatusOK, item)
}

// Delete deletes an item
// @Summary Delete an item
// @Description Delete an existing item by its ID
// @Tags items
// @Accept json
// @Produce json
// @Param id path int true "Item ID"
// @Success 204
// @Failure 400 {object} utils.ErrorResponse
// @Router /api/items/{id} [delete]
func (h *ItemHandler) Delete(w http.ResponseWriter, r *http.Request) {
	idStr := r.Context().Value(idKey).(string)
	id, err := strconv.Atoi(idStr)
	if err != nil {
		_ = utils.WriteError(w, utils.ErrorBadRequest("invalid item ID"))
		return
	}

	if err := h.repo.Delete(id); err != nil {
		h.logger.Error("failed to delete item", slog.String("error", err.Error()))
		_ = utils.WriteError(w, utils.ErrorInternalServer("failed to delete item"))
		return
	}

	utils.WriteNoContent(w)
}

// decodeJSON safely decodes JSON from a reader with size limit
func decodeJSON(reader io.ReadCloser, v interface{}) error {
	defer reader.Close()

	decoder := utils.NewLimitedDecoder(reader)
	if err := decoder.Decode(v); err != nil {
		if errors.Is(err, utils.ErrBodyTooLarge) {
			return utils.ErrBodyTooLarge
		}
		return err
	}

	return nil
}

// SetIDInContext returns a middleware that adds the item ID to the request context
func SetIDInContext(id string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx := context.WithValue(r.Context(), idKey, id)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
