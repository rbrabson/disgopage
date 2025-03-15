package disgopage

import (
	"fmt"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// Paginator represents a paginator that may be used to create paginated messages on a Discord server.
type Paginator struct {
	id       string
	config   *config
	messages map[string]*message
	mutex    sync.Mutex
	manager  *paginatorManager
}

// NewPaginator creates a new paginator.
func NewPaginator(opts ...ConfigOpt) *Paginator {
	config := defaultConfig
	config.Apply(opts)
	p := &Paginator{
		id:       fmt.Sprintf("paginator-%d", time.Now().UnixNano()),
		config:   &config,
		messages: make(map[string]*message),
		mutex:    sync.Mutex{},
		manager:  manager,
	}

	p.manager.addPaginator(p)

	log.WithFields(log.Fields{"itemsPerPage": p.config.ItemsPerPage, "idleWait": p.config.IdleWait}).Debug("created new paginator")
	return p
}

// createInteractionResponse creates and sends a message with the paginator's content.
func (p *Paginator) CreateInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate, title string, embedFields []*discordgo.MessageEmbedField, ephemeral ...bool) error {
	m := newMessage(p, title, embedFields)
	m.id = fmt.Sprintf("%s-%d", i.ChannelID, time.Now().UnixNano())
	m.interaction = i.Interaction
	m.ephemeral = len(ephemeral) > 0 && ephemeral[0]
	var flags discordgo.MessageFlags
	if m.ephemeral {
		flags = discordgo.MessageFlagsEphemeral
	}
	p.mutex.Lock()
	p.messages[m.id] = m
	p.mutex.Unlock()

	embeds := []*discordgo.MessageEmbed{m.makeEmbed()}
	components := []discordgo.MessageComponent{m.makeComponent(false)}
	m.registerComponentHandlers()
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     embeds,
			Components: components,
			Flags:      flags,
		},
	})
	if err != nil {
		log.WithFields(log.Fields{"paginator": p.id, "message": m.id, "channel": i.ChannelID, "error": err}).Error("error sending paginated message")
		m.deregisterComponentHandlers()
		p.mutex.Lock()
		delete(p.messages, m.id)
		p.mutex.Unlock()
		return err
	}
	log.WithFields(log.Fields{"paginator": p.id, "message": m.id, "channel": i.ChannelID}).Debug("created paginated message")
	return nil
}

// createInteractionResponse creates and sends a message with the paginator's content.
func (p *Paginator) CreateMessage(s *discordgo.Session, channelID string, title string, embedFields []*discordgo.MessageEmbedField) error {
	m := newMessage(p, title, embedFields)
	m.id = fmt.Sprintf("%s-%d", channelID, time.Now().UnixNano())
	m.channelID = channelID
	p.mutex.Lock()
	p.messages[m.id] = m
	p.mutex.Unlock()

	embeds := []*discordgo.MessageEmbed{m.makeEmbed()}
	components := []discordgo.MessageComponent{m.makeComponent(false)}
	m.registerComponentHandlers()

	message, err := s.ChannelMessageSendComplex(m.channelID, &discordgo.MessageSend{
		Embeds:     embeds,
		Components: components,
	})
	m.messageID = message.ID
	if err != nil {
		log.WithFields(log.Fields{"paginator": p.id, "message": m.id, "channel": channelID, "error": err}).Error("error sending paginated message")
		m.deregisterComponentHandlers()
		p.mutex.Lock()
		delete(p.messages, m.id)
		p.mutex.Unlock()
		return err
	}
	log.WithFields(log.Fields{"paginator": p.id, "message": m.id, "channel": channelID}).Debug("created paginated message")
	return nil
}

// Close closes the paginator and disables all paginated messages
func (p *Paginator) Close() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, m := range p.messages {
		m.disable()
		m.deregisterComponentHandlers()
		delete(p.messages, m.id)
	}

	manager.removePaginator(p)
}

// cleanup cleans up expired paginated messages. It is called by the manager's cleanup goroutine.
func (p *Paginator) cleanup() {
	p.mutex.Lock()
	defer p.mutex.Unlock()

	for _, m := range p.messages {
		if m.hasExpired() {
			m.disable()
			m.deregisterComponentHandlers()
			delete(p.messages, m.id)
		}
	}
}
