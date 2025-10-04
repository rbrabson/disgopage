package disgopage

import (
	"testing"

	"github.com/bwmarrin/discordgo"
)

func TestNewMessage(t *testing.T) {
	// Create a paginator
	p := &Paginator{
		id:       "test-paginator",
		config:   &defaultConfig,
		messages: make(map[string]*message),
	}

	// Create test embed fields
	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:  "Field 1",
			Value: "Value 1",
		},
		{
			Name:  "Field 2",
			Value: "Value 2",
		},
	}

	// Create a new message
	title := "Test Message"
	msg := newMessage(p, title, embedFields)

	// Verify the message was created correctly
	if msg == nil {
		t.Errorf("Expected message to be created, got nil")
		return
	}
	if msg.paginator != p {
		t.Errorf("Expected message.paginator to be set correctly")
	}
	if msg.title != title {
		t.Errorf("Expected message.title to be %s, got %s", title, msg.title)
	}
	if len(msg.embedFields) != len(embedFields) {
		t.Errorf("Expected message to have %d embed fields, got %d", len(embedFields), len(msg.embedFields))
	}
	if msg.currentPage != 0 {
		t.Errorf("Expected message.currentPage to be 0, got %d", msg.currentPage)
	}
	if msg.expiry.IsZero() {
		t.Errorf("Expected message.expiry to be set")
	}
}

func TestPageCount(t *testing.T) {
	// Create a paginator with 5 items per page
	p := &Paginator{
		id: "test-paginator",
		config: &config{
			ItemsPerPage: 5,
		},
		messages: make(map[string]*message),
	}

	// Test cases with different numbers of embed fields
	testCases := []struct {
		name          string
		embedFields   []*discordgo.MessageEmbedField
		expectedPages int
	}{
		{
			name:          "No fields",
			embedFields:   []*discordgo.MessageEmbedField{},
			expectedPages: 1, // Should have at least 1 page even with no fields
		},
		{
			name:          "Exactly one page",
			embedFields:   make([]*discordgo.MessageEmbedField, 5),
			expectedPages: 1,
		},
		{
			name:          "Partial second page",
			embedFields:   make([]*discordgo.MessageEmbedField, 7),
			expectedPages: 2,
		},
		{
			name:          "Multiple full pages",
			embedFields:   make([]*discordgo.MessageEmbedField, 15),
			expectedPages: 3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// Create a message with the test case's embed fields
			msg := newMessage(p, "Test", tc.embedFields)

			// Check the page count
			pageCount := msg.pageCount()
			if pageCount != tc.expectedPages {
				t.Errorf("Expected %d pages, got %d", tc.expectedPages, pageCount)
			}
		})
	}
}

func TestMakeEmbed(t *testing.T) {
	// Create a paginator with custom embed color
	p := &Paginator{
		id: "test-paginator",
		config: &config{
			EmbedColor:   0xFF0000,
			ItemsPerPage: 2,
		},
		messages: make(map[string]*message),
	}

	// Create test embed fields
	embedFields := []*discordgo.MessageEmbedField{
		{
			Name:  "Field 1",
			Value: "Value 1",
		},
		{
			Name:  "Field 2",
			Value: "Value 2",
		},
		{
			Name:  "Field 3",
			Value: "Value 3",
		},
	}

	// Create a message
	title := "Test Embed"
	msg := newMessage(p, title, embedFields)

	// Test the embed for page 0 (first page)
	embed := msg.makeEmbed()

	// Verify the embed properties
	if embed == nil {
		t.Errorf("Expected embed to be created, got nil")
		return
	}
	if embed.Title != title {
		t.Errorf("Expected embed.Title to be %s, got %s", title, embed.Title)
	}
	if embed.Color != 0xFF0000 {
		t.Errorf("Expected embed.Color to be 0xFF0000, got %d", embed.Color)
	}

	// Verify the embed has the correct fields for the first page
	if len(embed.Fields) != 2 {
		t.Errorf("Expected embed to have 2 fields on first page, got %d", len(embed.Fields))
	}
	if embed.Fields[0].Name != "Field 1" {
		t.Errorf("Expected first field name to be Field 1, got %s", embed.Fields[0].Name)
	}
	if embed.Fields[1].Name != "Field 2" {
		t.Errorf("Expected second field name to be Field 2, got %s", embed.Fields[1].Name)
	}

	// Test pagination - move to page 1 and check fields
	msg.currentPage = 1
	embed = msg.makeEmbed()

	if len(embed.Fields) != 1 {
		t.Errorf("Expected embed to have 1 field on second page, got %d", len(embed.Fields))
	}
	if embed.Fields[0].Name != "Field 3" {
		t.Errorf("Expected first field on second page to be Field 3, got %s", embed.Fields[0].Name)
	}
}

func TestCustomButtonID(t *testing.T) {
	// Create a paginator with a custom ID prefix
	p := &Paginator{
		id: "test-paginator",
		config: &config{
			CustomIDPrefix: "custom-prefix",
		},
		messages: make(map[string]*message),
	}

	// Create a message
	msg := &message{
		id:        "test-message",
		paginator: p,
	}

	// Test the custom button ID
	buttonID := msg.customButtonID("next")
	expectedID := "custom-prefix:test-paginator:test-message:next"

	if buttonID != expectedID {
		t.Errorf("Expected button ID to be %s, got %s", expectedID, buttonID)
	}
}
