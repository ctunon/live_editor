// handlers/document.go
package handlers

import (
	"encoding/json"
	"live_editor/models"
	"live_editor/storage"
	"live_editor/utils"
	"net/http"
)

// DocumentHandler provides HTTP handlers for document operations
type DocumentHandler struct {
	storage *storage.MemoryStorage
}

// NewDocumentHandler creates a new DocumentHandler
func NewDocumentHandler(storage *storage.MemoryStorage) *DocumentHandler {
	return &DocumentHandler{storage: storage}
}

func (h *DocumentHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		h.GetDocument(w, r)
	case http.MethodPost:
		h.CreateDocument(w, r)
	case http.MethodPut:
		h.UpdateDocument(w, r)
	case http.MethodDelete:
		h.DeleteDocument(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// CreateDocument godoc
// @Summary Create a new document
// @Description Create a new document with content
// @Tags documents
// @Accept  json
// @Produce  json
// @Param document body models.Document true "Document content"
// @Success 201 {object} models.Document
// @Failure 400 {object} map[string]string
// @Router /documents [post]
func (h *DocumentHandler) CreateDocument(w http.ResponseWriter, r *http.Request) {
	var doc models.Document
	err := json.NewDecoder(r.Body).Decode(&doc)
	if err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}
	h.storage.CreateDocument(doc)
	utils.SendJSON(w, doc, http.StatusCreated)

}

// GetDocument godoc
// @Summary Get a document by ID
// @Description Get a document by its ID
// @Tags documents
// @Produce  json
// @Param id query string true "Document ID"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /documents [get]
func (h *DocumentHandler) GetDocument(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.SendError(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	doc, exists := h.storage.GetDocument(id)
	if !exists {
		utils.SendError(w, "Document not found", http.StatusNotFound)
		return
	}

	utils.SendJSON(w, doc, http.StatusOK)
}

// UpdateDocument godoc
// @Summary Update a document
// @Description Update the content of a document by ID
// @Tags documents
// @Accept  json
// @Produce  json
// @Param id query string true "Document ID"
// @Param document body models.Document true "Updated document content"
// @Success 200 {object} models.Document
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /documents [put]
func (h *DocumentHandler) UpdateDocument(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.SendError(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	var updatedDoc models.Document
	err := json.NewDecoder(r.Body).Decode(&updatedDoc)
	if err != nil {
		utils.SendError(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	updatedDoc.ID = id // ensure the ID is not changed
	if !h.storage.UpdateDocument(updatedDoc) {
		utils.SendError(w, "Document not found", http.StatusNotFound)
		return
	}

	utils.SendJSON(w, updatedDoc, http.StatusOK)
}

// DeleteDocument godoc
// @Summary Delete a document
// @Description Delete a document by ID
// @Tags documents
// @Param id query string true "Document ID"
// @Success 204
// @Failure 400 {object} map[string]string
// @Failure 404 {object} map[string]string
// @Router /documents [delete]
func (h *DocumentHandler) DeleteDocument(w http.ResponseWriter, r *http.Request) {
	id := r.URL.Query().Get("id")
	if id == "" {
		utils.SendError(w, "Document ID is required", http.StatusBadRequest)
		return
	}

	if !h.storage.DeleteDocument(id) {
		utils.SendError(w, "Document not found", http.StatusNotFound)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}
