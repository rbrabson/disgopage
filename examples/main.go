package main

import (
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/joho/godotenv"
	page "github.com/rbrabson/disgopage"
)

const (
	botIntents = discordgo.IntentGuilds |
		discordgo.IntentGuildMessages |
		discordgo.IntentDirectMessages |
		discordgo.IntentGuildEmojis
)

var (
	Token string
	AppID string
)

var (
	dg *discordgo.Session
)

var (
	componentHandlers = make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate))
	commandHandlers   = make(map[string]func(*discordgo.Session, *discordgo.InteractionCreate))
	commands          = []*discordgo.ApplicationCommand{
		{
			Name:        "paginator",
			Description: "Paginator Command",
		},
	}
)

var (
	embeds = []*discordgo.MessageEmbedField{
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
		{
			Name:  "Field 4",
			Value: "Value 4",
		},
		{
			Name:  "Field 5",
			Value: "Value 5",
		},
		{
			Name:  "Field 6",
			Value: "Value 6",
		},
		{
			Name:  "Field 7",
			Value: "Value 7",
		},
	}
)

func init() {
	if err := godotenv.Load(".env_test"); err != nil {
		slog.Error("Error loading .env file",
			slog.Any("error", err),
		)
	}
	os.Getenv("DISCORD_BOT_TOKEN")
	AppID = os.Getenv("DISCORD_APP_ID")
}

func main() {
	var err error
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		slog.Error("error creating Discord session,",
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	dg.Identify.Intents = botIntents

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		slog.Info("Bot is up!")
	})

	commandHandlers["paginator"] = paginator
	dg.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand:
			if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
				h(s, i)
			}
		case discordgo.InteractionMessageComponent:
			if h, ok := componentHandlers[i.MessageComponentData().CustomID]; ok {
				h(s, i)
			}
		}
	})

	_, err = dg.ApplicationCommandBulkOverwrite(AppID, "", commands)
	if err != nil {
		slog.Error("failed to load bot commands",
			slog.String("appID", AppID),
			slog.Any("commands", commands),
			slog.Any("error", err),
		)
		os.Exit(1)
	}

	if err := dg.Open(); err != nil {
		slog.Error("error opening connection,",
			slog.Any("error", err),
		)
		os.Exit(1)
	}
	defer func() {
		if err := dg.Close(); err != nil {
			slog.Error("error closing Discord session",
				slog.Any("error", err),
			)
		}
	}()

	p := page.NewPaginator(
		page.WithDiscordConfig(
			page.DiscordConfig{
				Session:                dg,
				AddComponentHandler:    addComponentHandler,
				RemoveComponentHandler: removeComponentHandler,
			},
		),
	)
	if err := p.CreateMessage(dg, "1135713066164703232", "Paginator Using CreateMessage", embeds); err != nil {
		slog.Error("error creating message",
			slog.Any("error", err),
		)
	}

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	slog.Info("Press Ctrl+C to exit")
	<-sc
}

// process the `/paginator` command
func paginator(s *discordgo.Session, i *discordgo.InteractionCreate) {
	p := page.NewPaginator(
		page.WithDiscordConfig(
			page.DiscordConfig{
				Session:                dg,
				AddComponentHandler:    addComponentHandler,
				RemoveComponentHandler: removeComponentHandler,
			},
		),
	)
	if err := p.CreateInteractionResponse(s, i, "Paginator Using CreateInteractionResponse", embeds, true); err != nil {
		slog.Error("error creating interaction response",
			slog.Any("error", err),
		)
	}
}

func addComponentHandler(key string, handler func(*discordgo.Session, *discordgo.InteractionCreate)) {
	componentHandlers[key] = handler
}

func removeComponentHandler(key string) {
	delete(componentHandlers, key)
}
