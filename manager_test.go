package disgopage

import (
	"testing"
)

func TestNewManager(t *testing.T) {
	// Create a new manager
	m := newManager()

	// Verify the manager was created correctly
	if m == nil {
		t.Errorf("Expected manager to be created, got nil")
		return
	}
	if m.paginators == nil {
		t.Errorf("Expected paginators map to be initialized")
	}
	if len(m.paginators) != 0 {
		t.Errorf("Expected paginators map to be empty, got %d items", len(m.paginators))
	}
}

func TestAddPaginator(t *testing.T) {
	// Create a new manager and paginator
	m := newManager()
	p := &Paginator{
		id:       "test-paginator",
		config:   &defaultConfig,
		messages: make(map[string]*message),
	}

	// Add the paginator to the manager
	m.addPaginator(p)

	// Verify the paginator was added
	if len(m.paginators) != 1 {
		t.Errorf("Expected paginators map to have 1 item, got %d", len(m.paginators))
	}
	if m.paginators["test-paginator"] != p {
		t.Errorf("Expected paginator to be added with correct ID")
	}
}

func TestRemovePaginator(t *testing.T) {
	// Create a new manager and paginator
	m := newManager()
	p := &Paginator{
		id:       "test-paginator",
		config:   &defaultConfig,
		messages: make(map[string]*message),
	}

	// Add the paginator to the manager
	m.paginators["test-paginator"] = p

	// Remove the paginator
	m.removePaginator(p)

	// Verify the paginator was removed
	if len(m.paginators) != 0 {
		t.Errorf("Expected paginators map to be empty after removal, got %d items", len(m.paginators))
	}
	if _, exists := m.paginators["test-paginator"]; exists {
		t.Errorf("Expected paginator to be removed")
	}
}

func TestCleanup(t *testing.T) {
	// Since we can't directly mock the Paginator.cleanup method,
	// we'll just test that the manager's cleanup method doesn't panic

	// Create a new manager
	m := newManager()

	// Create a paginator
	p := &Paginator{
		id:       "test-paginator",
		config:   &defaultConfig,
		messages: make(map[string]*message),
	}

	// Add the paginator to the manager
	m.paginators["test-paginator"] = p

	// Call cleanup - this should not panic
	m.cleanup()

	// Verify the paginator is still in the manager (cleanup doesn't remove it)
	if _, exists := m.paginators["test-paginator"]; !exists {
		t.Errorf("Expected paginator to still exist after cleanup")
	}
}
