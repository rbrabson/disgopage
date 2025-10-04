package disgopage

import (
	"testing"
	"time"
)

func TestNewPaginator(t *testing.T) {
	// Test creating a paginator with default options
	p := NewPaginator()

	// Verify the paginator was created correctly
	if p == nil {
		t.Errorf("Expected paginator to be created, got nil")
		return
	}
	if p.id == "" {
		t.Errorf("Expected paginator.id to be set")
	}
	if p.config == nil {
		t.Errorf("Expected paginator.config to be set")
	}
	if p.messages == nil {
		t.Errorf("Expected paginator.messages to be initialized")
	}
	if len(p.messages) != 0 {
		t.Errorf("Expected paginator.messages to be empty, got %d items", len(p.messages))
	}
	if p.manager == nil {
		t.Errorf("Expected paginator.manager to be set")
	}

	// Verify default config values
	if p.config.ItemsPerPage != defaultConfig.ItemsPerPage {
		t.Errorf("Expected ItemsPerPage to be %d, got %d", defaultConfig.ItemsPerPage, p.config.ItemsPerPage)
	}
	if p.config.IdleWait != defaultConfig.IdleWait {
		t.Errorf("Expected IdleWait to be %s, got %s", defaultConfig.IdleWait, p.config.IdleWait)
	}
	if p.config.CustomIDPrefix != defaultConfig.CustomIDPrefix {
		t.Errorf("Expected CustomIDPrefix to be %s, got %s", defaultConfig.CustomIDPrefix, p.config.CustomIDPrefix)
	}
}

func TestNewPaginatorWithOptions(t *testing.T) {
	// Test creating a paginator with custom options
	customPrefix := "custom-prefix"
	customItemsPerPage := 10
	customIdleWait := time.Minute * 10

	p := NewPaginator(
		WithCustomIDPrefix(customPrefix),
		WithItemsPerPage(customItemsPerPage),
		WithIdleWait(customIdleWait),
	)

	// Verify the custom options were applied
	if p.config.CustomIDPrefix != customPrefix {
		t.Errorf("Expected CustomIDPrefix to be %s, got %s", customPrefix, p.config.CustomIDPrefix)
	}
	if p.config.ItemsPerPage != customItemsPerPage {
		t.Errorf("Expected ItemsPerPage to be %d, got %d", customItemsPerPage, p.config.ItemsPerPage)
	}
	if p.config.IdleWait != customIdleWait {
		t.Errorf("Expected IdleWait to be %s, got %s", customIdleWait, p.config.IdleWait)
	}
}

func TestClose(t *testing.T) {
	// Create a paginator
	p := NewPaginator()

	// Add a mock message
	mockMessage := &message{
		id:        "test-message",
		paginator: p,
		expiry:    time.Now().Add(time.Hour),
	}
	p.messages["test-message"] = mockMessage

	// Close the paginator
	p.Close()

	// Verify the messages map is empty after closing
	if len(p.messages) != 0 {
		t.Errorf("Expected messages map to be empty after Close(), got %d items", len(p.messages))
	}
}

func TestPaginatorCleanup(t *testing.T) {
	// Create a paginator
	p := NewPaginator()

	// Add a mock message that has expired
	expiredMessage := &message{
		id:        "expired-message",
		paginator: p,
		expiry:    time.Now().Add(-time.Hour), // Expired 1 hour ago
	}
	p.messages["expired-message"] = expiredMessage

	// Add a mock message that has not expired
	activeMessage := &message{
		id:        "active-message",
		paginator: p,
		expiry:    time.Now().Add(time.Hour), // Expires in 1 hour
	}
	p.messages["active-message"] = activeMessage

	// Call cleanup
	p.cleanup()

	// Verify the expired message was removed
	if _, exists := p.messages["expired-message"]; exists {
		t.Errorf("Expected expired message to be removed")
	}

	// Verify the active message still exists
	if _, exists := p.messages["active-message"]; !exists {
		t.Errorf("Expected active message to still exist")
	}
}
