package disgopage

import (
	"time"

	"github.com/bwmarrin/discordgo"
)

// defaultConfig is the default configuration used by the paginator.
var defaultConfig = config{
	ButtonsConfig: ButtonsConfig{
		First: &ComponentOption{
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏮️",
			},
			Style: discordgo.PrimaryButton,
		},
		Back: &ComponentOption{
			Emoji: &discordgo.ComponentEmoji{
				Name: "◀️",
			},
			Style: discordgo.PrimaryButton,
		},
		Next: &ComponentOption{
			Emoji: &discordgo.ComponentEmoji{
				Name: "▶️",
			},
			Style: discordgo.PrimaryButton,
		},
		Last: &ComponentOption{
			Emoji: &discordgo.ComponentEmoji{
				Name: "⏭️",
			},
			Style: discordgo.PrimaryButton,
		},
	},
	CustomIDPrefix: "paginator",
	EmbedColor:     0x4c50c1,
	ItemsPerPage:   5,
	IdleWait:       time.Minute * 5,
}

// config is the configuration used by the paginator.
type config struct {
	ButtonsConfig  ButtonsConfig
	CustomIDPrefix string
	EmbedColor     int
	ItemsPerPage   int
	DiscordConfig  DiscordConfig
	IdleWait       time.Duration
}

// ComponentOption are the options used to create a pagination button.
type ComponentOption struct {
	Emoji *discordgo.ComponentEmoji
	Label string
	Style discordgo.ButtonStyle
}

// ConfigOpt is a function that can be used to modify the paginator's configuration.
type ConfigOpt func(config *config)

// Apply applies the given RequestOpt(s) to the RequestConfig & sets the context if none is set
func (c *config) Apply(opts []ConfigOpt) {
	for _, opt := range opts {
		opt(c)
	}
}

// ButtonsConfig is the configuration for the pagination buttons.
type ButtonsConfig struct {
	First *ComponentOption
	Back  *ComponentOption
	Stop  *ComponentOption
	Next  *ComponentOption
	Last  *ComponentOption
}

// DiscordConfig is the configuration used by the paginator when using Discord.
type DiscordConfig struct {
	Session                *discordgo.Session
	AddComponentHandler    func(key string, handler func(*discordgo.Session, *discordgo.InteractionCreate))
	RemoveComponentHandler func(key string)
}

// WithButtonsConfig sets the button configuration for the paginator.
func WithButtonsConfig(buttonsConfig ButtonsConfig) ConfigOpt {
	return func(config *config) {
		config.ButtonsConfig = buttonsConfig
	}
}

// WithCustomIDPrefix sets the custom ID prefix for the paginator.
func WithCustomIDPrefix(prefix string) ConfigOpt {
	return func(config *config) {
		config.CustomIDPrefix = prefix
	}
}

// WithEmbedColor sets the embed color for the paginator.
func WithEmbedColor(color int) ConfigOpt {
	return func(config *config) {
		config.EmbedColor = color
	}
}

// WithDiscordConfig sets the Discord configuration for the paginator.
func WithDiscordConfig(discordConfig DiscordConfig) ConfigOpt {
	return func(config *config) {
		config.DiscordConfig = discordConfig
	}
}

// WithItemsPerPage sets the default number of items per page for the paginator.
func WithItemsPerPage(itemsPerPage int) ConfigOpt {
	return func(config *config) {
		config.ItemsPerPage = itemsPerPage
	}
}

// WithIdleWait sets the idle wait time for the paginator.
func WithIdleWait(idleWait time.Duration) ConfigOpt {
	return func(config *config) {
		config.IdleWait = idleWait
	}
}
