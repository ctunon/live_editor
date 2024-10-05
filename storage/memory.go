// storage/memory.go
package storage

import (
	"live_editor/models"
	"sync"
)

// MemoryStorage manages in-memory document storage
type MemoryStorage struct {
	documents map[string]models.Document
	mu        sync.Mutex
}

// NewMemoryStorage creates a new in-memory storage instance
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		documents: make(map[string]models.Document),
	}
}

// CreateDocument adds a new document to the storage
func (s *MemoryStorage) CreateDocument(doc models.Document) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.documents[doc.ID] = doc
}

// GetDocument retrieves a document by its ID
func (s *MemoryStorage) GetDocument(id string) (models.Document, bool) {
	s.mu.Lock()
	defer s.mu.Unlock()
	doc, exists := s.documents[id]
	return doc, exists
}

// UpdateDocument updates an existing document in the storage
func (s *MemoryStorage) UpdateDocument(doc models.Document) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.documents[doc.ID]; !exists {
		return false
	}
	s.documents[doc.ID] = doc
	return true
}

// DeleteDocument removes a document by its ID
func (s *MemoryStorage) DeleteDocument(id string) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.documents[id]; !exists {
		return false
	}
	delete(s.documents, id)
	return true
}
