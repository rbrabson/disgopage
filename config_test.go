package disgopage

import (
	"testing"
	"time"

	"github.com/bwmarrin/discordgo"
)

func TestConfigDefaults(t *testing.T) {
	// Create a new config with default values
	cfg := defaultConfig

	// Verify default values
	if cfg.ButtonsConfig.First.Emoji.Name != "⏮️" {
		t.Errorf("Expected First button emoji to be ⏮️, got %s", cfg.ButtonsConfig.First.Emoji.Name)
	}
	if cfg.ButtonsConfig.Back.Emoji.Name != "◀️" {
		t.Errorf("Expected Back button emoji to be ◀️, got %s", cfg.ButtonsConfig.Back.Emoji.Name)
	}
	if cfg.ButtonsConfig.Next.Emoji.Name != "▶️" {
		t.Errorf("Expected Next button emoji to be ▶️, got %s", cfg.ButtonsConfig.Next.Emoji.Name)
	}
	if cfg.ButtonsConfig.Last.Emoji.Name != "⏭️" {
		t.Errorf("Expected Last button emoji to be ⏭️, got %s", cfg.ButtonsConfig.Last.Emoji.Name)
	}
	if cfg.CustomIDPrefix != "paginator" {
		t.Errorf("Expected CustomIDPrefix to be paginator, got %s", cfg.CustomIDPrefix)
	}
	if cfg.EmbedColor != 0x4c50c1 {
		t.Errorf("Expected EmbedColor to be 0x4c50c1, got %d", cfg.EmbedColor)
	}
	if cfg.ItemsPerPage != 5 {
		t.Errorf("Expected ItemsPerPage to be 5, got %d", cfg.ItemsPerPage)
	}
	if cfg.IdleWait != time.Minute*5 {
		t.Errorf("Expected IdleWait to be 5 minutes, got %s", cfg.IdleWait)
	}
}

func TestWithButtonsConfig(t *testing.T) {
	// Create a custom buttons config
	customButtons := ButtonsConfig{
		First: &ComponentOption{
			Emoji: &discordgo.ComponentEmoji{
				Name: "🔍",
			},
			Style: discordgo.SecondaryButton,
		},
		Back: &ComponentOption{
			Emoji: &discordgo.ComponentEmoji{
				Name: "⬅️",
			},
			Style: discordgo.SecondaryButton,
		},
	}

	// Create a config with the custom buttons
	cfg := defaultConfig
	opt := WithButtonsConfig(customButtons)
	opt(&cfg)

	// Verify the buttons were updated
	if cfg.ButtonsConfig.First.Emoji.Name != "🔍" {
		t.Errorf("Expected First button emoji to be 🔍, got %s", cfg.ButtonsConfig.First.Emoji.Name)
	}
	if cfg.ButtonsConfig.First.Style != discordgo.SecondaryButton {
		t.Errorf("Expected First button style to be SecondaryButton, got %v", cfg.ButtonsConfig.First.Style)
	}
	if cfg.ButtonsConfig.Back.Emoji.Name != "⬅️" {
		t.Errorf("Expected Back button emoji to be ⬅️, got %s", cfg.ButtonsConfig.Back.Emoji.Name)
	}
	if cfg.ButtonsConfig.Back.Style != discordgo.SecondaryButton {
		t.Errorf("Expected Back button style to be SecondaryButton, got %v", cfg.ButtonsConfig.Back.Style)
	}

	// Verify other buttons remain unchanged from default
	if cfg.ButtonsConfig.Next.Emoji.Name != "▶️" {
		t.Errorf("Expected Next button emoji to be ▶️, got %s", cfg.ButtonsConfig.Next.Emoji.Name)
	}
	if cfg.ButtonsConfig.Next.Style != discordgo.PrimaryButton {
		t.Errorf("Expected Next button style to be PrimaryButton, got %v", cfg.ButtonsConfig.Next.Style)
	}
}

func TestWithCustomIDPrefix(t *testing.T) {
	// Create a config with a custom prefix
	cfg := defaultConfig
	opt := WithCustomIDPrefix("custom-prefix")
	opt(&cfg)

	// Verify the prefix was updated
	if cfg.CustomIDPrefix != "custom-prefix" {
		t.Errorf("Expected CustomIDPrefix to be custom-prefix, got %s", cfg.CustomIDPrefix)
	}
}

func TestWithEmbedColor(t *testing.T) {
	// Create a config with a custom embed color
	cfg := defaultConfig
	opt := WithEmbedColor(0xFF0000)
	opt(&cfg)

	// Verify the color was updated
	if cfg.EmbedColor != 0xFF0000 {
		t.Errorf("Expected EmbedColor to be 0xFF0000, got %d", cfg.EmbedColor)
	}
}

func TestWithItemsPerPage(t *testing.T) {
	// Create a config with a custom items per page
	cfg := defaultConfig
	opt := WithItemsPerPage(10)
	opt(&cfg)

	// Verify the items per page was updated
	if cfg.ItemsPerPage != 10 {
		t.Errorf("Expected ItemsPerPage to be 10, got %d", cfg.ItemsPerPage)
	}
}

func TestWithIdleWait(t *testing.T) {
	// Create a config with a custom idle wait time
	cfg := defaultConfig
	opt := WithIdleWait(time.Minute * 10)
	opt(&cfg)

	// Verify the idle wait time was updated
	if cfg.IdleWait != time.Minute*10 {
		t.Errorf("Expected IdleWait to be 10 minutes, got %s", cfg.IdleWait)
	}
}

func TestWithDiscordConfig(t *testing.T) {
	// Create a mock session and handlers
	session := &discordgo.Session{}
	addHandler := func(key string, handler func(*discordgo.Session, *discordgo.InteractionCreate)) {}
	removeHandler := func(key string) {}

	// Create a config with the discord config
	cfg := defaultConfig
	discordCfg := DiscordConfig{
		Session:                session,
		AddComponentHandler:    addHandler,
		RemoveComponentHandler: removeHandler,
	}
	opt := WithDiscordConfig(discordCfg)
	opt(&cfg)

	// Verify the discord config was updated
	if cfg.DiscordConfig.Session != session {
		t.Errorf("Expected Session to be set correctly")
	}
	if cfg.DiscordConfig.AddComponentHandler == nil {
		t.Errorf("Expected AddComponentHandler to be set")
	}
	if cfg.DiscordConfig.RemoveComponentHandler == nil {
		t.Errorf("Expected RemoveComponentHandler to be set")
	}
}

func TestConfigApply(t *testing.T) {
	// Create a config with multiple options
	cfg := defaultConfig
	opts := []ConfigOpt{
		WithCustomIDPrefix("test-prefix"),
		WithEmbedColor(0x00FF00),
		WithItemsPerPage(15),
	}

	// Apply the options
	cfg.Apply(opts)

	// Verify all options were applied
	if cfg.CustomIDPrefix != "test-prefix" {
		t.Errorf("Expected CustomIDPrefix to be test-prefix, got %s", cfg.CustomIDPrefix)
	}
	if cfg.EmbedColor != 0x00FF00 {
		t.Errorf("Expected EmbedColor to be 0x00FF00, got %d", cfg.EmbedColor)
	}
	if cfg.ItemsPerPage != 15 {
		t.Errorf("Expected ItemsPerPage to be 15, got %d", cfg.ItemsPerPage)
	}
}
