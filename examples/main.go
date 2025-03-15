package main

import (
	"os"
	"os/signal"
	"syscall"

	log "github.com/sirupsen/logrus"

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
	godotenv.Load(".env_test")
	Token = os.Getenv("DISCORD_BOT_TOKEN")
	AppID = os.Getenv("DISCORD_APP_ID")
}

func main() {
	var err error
	dg, err = discordgo.New("Bot " + Token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
		return
	}

	dg.Identify.Intents = botIntents

	dg.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Info("Bot is up!")
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
		log.WithFields(log.Fields{"appID": AppID, "commands": commands, "error": err}).Fatal("failed to load bot commands")
	}

	dg.Open()
	defer dg.Close()

	p := page.NewPaginator(
		page.WithDiscordConfig(
			page.DiscordConfig{
				Session:                dg,
				AddComponentHandler:    addComponentHandler,
				RemoveComponentHandler: removeComponentHandler,
			},
		),
	)
	p.CreateMessage(dg, "1135713066164703232", "Paginator Using CreateMessage", embeds)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	log.Info("Press Ctrl+C to exit")
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
	p.CreateInteractionResponse(s, i, "Paginator Using CreateInteractionResponse", embeds, true)
}

func addComponentHandler(key string, handler func(*discordgo.Session, *discordgo.InteractionCreate)) {
	componentHandlers[key] = handler
}

func removeComponentHandler(key string) {
	delete(componentHandlers, key)
}
