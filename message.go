package disgopage

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// message represents a single message in the paginator. It contains the data
// to be paginated, as well as the state of the paginator.
type message struct {
	id          string
	title       string
	embedFields []*discordgo.MessageEmbedField
	expiry      time.Time
	currentPage int
	channelID   string
	paginator   *Paginator
	interaction *discordgo.Interaction
	ephemeral   bool
}

// newMessge creates a new message for the paginator.
func newInteractionResponse(p *Paginator, title string, embedFields []*discordgo.MessageEmbedField) *message {
	return &message{
		paginator:   p,
		title:       title,
		embedFields: embedFields,
		expiry:      time.Now().Add(p.config.IdleWait),
	}
}

// editMessage edits the current message sent by the paginator in a channel.
func (m *message) editInteractionResponse(s *discordgo.Session, i *discordgo.InteractionCreate) error {
	embeds := []*discordgo.MessageEmbed{m.makeEmbed()}
	components := []discordgo.MessageComponent{m.makeComponent(false)}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseDeferredMessageUpdate,
	})

	_, err := s.InteractionResponseEdit(m.interaction, &discordgo.WebhookEdit{
		Embeds:     &embeds,
		Components: &components,
	})
	if err != nil {
		log.WithFields(log.Fields{"paginator": m.id, "channel": m.channelID, "error": err}).Error("error editing paginated message")
		return err
	}

	log.WithFields(log.Fields{"paginator": m.id, "channel": m.channelID}).Debug("edited paginated message")
	return nil
}

// disable disables the message by removing the buttons and setting the setting the expiry time to now.
func (m *message) disable() error {
	embeds := []*discordgo.MessageEmbed{m.makeEmbed()}
	components := []discordgo.MessageComponent{m.makeComponent(true)}

	session := m.paginator.config.DiscordConfig.Session
	_, err := session.InteractionResponseEdit(m.interaction, &discordgo.WebhookEdit{
		Embeds:     &embeds,
		Components: &components,
	})
	if err != nil {
		log.WithFields(log.Fields{"paginator": m.id, "channel": m.channelID, "error": err}).Error("error disabling paginated message")
		return err
	}

	log.WithFields(log.Fields{"paginator": m.id, "channel": m.channelID}).Debug("disabled paginated message")
	return nil
}

// pageCount returns the number of pages in the paginator.
func (m *message) pageCount() int {
	itemsPerPage := m.getItemsPerPage()
	pageCount := (len(m.embedFields) + itemsPerPage - 1) / itemsPerPage
	return pageCount
}

// makeEmbed creates the message embed to be included for the current page.
func (m *message) makeEmbed() *discordgo.MessageEmbed {
	embed := &discordgo.MessageEmbed{
		Color:  m.paginator.config.EmbedColor,
		Title:  m.title,
		Fields: make([]*discordgo.MessageEmbedField, 0, m.getItemsPerPage()),
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page %d of %d", m.currentPage+1, m.pageCount()),
		},
	}
	start := m.currentPage * m.getItemsPerPage()
	end := min(start+m.getItemsPerPage(), len(m.embedFields))
	embed.Fields = append(embed.Fields, m.embedFields[start:end]...)
	return embed
}

// makeComponent creates  the message components to be included in the
// message. It returns an action row that contains the buttons used to navigate
// through the paginator.
func (m *message) makeComponent(disabled bool) discordgo.MessageComponent {
	cfg := m.paginator.config.ButtonsConfig
	actionRow := discordgo.ActionsRow{}

	if cfg.First != nil {
		buttonID := m.customButtonID("first")
		actionRow.Components = append(actionRow.Components, discordgo.Button{
			Label:    cfg.First.Label,
			Style:    cfg.First.Style,
			Disabled: disabled || m.currentPage == 0,
			Emoji:    cfg.First.Emoji,
			CustomID: buttonID,
		})
	}
	if cfg.Back != nil {
		buttonID := m.customButtonID("back")
		actionRow.Components = append(actionRow.Components, discordgo.Button{
			Label:    cfg.Back.Label,
			Style:    cfg.Back.Style,
			Disabled: disabled || m.currentPage == 0,
			Emoji:    cfg.Back.Emoji,
			CustomID: buttonID,
		})
	}
	if cfg.Stop != nil {
		buttonID := m.customButtonID("stop")
		actionRow.Components = append(actionRow.Components, discordgo.Button{
			Label:    cfg.Stop.Label,
			Style:    cfg.Stop.Style,
			Disabled: disabled || m.currentPage == 0,
			Emoji:    cfg.Stop.Emoji,
			CustomID: buttonID,
		})
	}
	if cfg.Next != nil {
		buttonID := m.customButtonID("next")
		actionRow.Components = append(actionRow.Components, discordgo.Button{
			Label:    cfg.Next.Label,
			Style:    cfg.Next.Style,
			Disabled: disabled || m.currentPage == m.pageCount()-1,
			Emoji:    cfg.Next.Emoji,
			CustomID: buttonID,
		})
	}
	if cfg.Last != nil {
		buttonID := m.customButtonID("last")
		actionRow.Components = append(actionRow.Components, discordgo.Button{
			Label:    cfg.Last.Label,
			Style:    cfg.Last.Style,
			Disabled: disabled || m.currentPage == m.pageCount()-1,
			Emoji:    cfg.Last.Emoji,
			CustomID: buttonID,
		})
	}

	return actionRow
}

// registerComponentHandlers registers the component handlers for the paginator.
func (m *message) registerComponentHandlers() {
	cfg := m.paginator.config
	if cfg.ButtonsConfig.First != nil {
		buttonID := m.customButtonID("first")
		cfg.DiscordConfig.AddComponentHandler(buttonID, pageThroughItems)
	}
	if cfg.ButtonsConfig.Back != nil {
		buttonID := m.customButtonID("back")
		cfg.DiscordConfig.AddComponentHandler(buttonID, pageThroughItems)
	}
	if cfg.ButtonsConfig.Stop != nil {
		buttonID := m.customButtonID("stop")
		cfg.DiscordConfig.AddComponentHandler(buttonID, pageThroughItems)
	}
	if cfg.ButtonsConfig.Next != nil {
		buttonID := m.customButtonID("next")
		cfg.DiscordConfig.AddComponentHandler(buttonID, pageThroughItems)
	}
	if cfg.ButtonsConfig.Last != nil {
		buttonID := m.customButtonID("last")
		cfg.DiscordConfig.AddComponentHandler(buttonID, pageThroughItems)
	}
	log.WithFields(log.Fields{"paginator": m.id}).Trace("registered component handlers")
}

// deregisterComponentHandlers deregisters the component handlers for the paginator.
func (m *message) deregisterComponentHandlers() {
	cfg := m.paginator.config
	if cfg.ButtonsConfig.First != nil {
		buttonID := m.customButtonID("first")
		cfg.DiscordConfig.RemoveComponentHandler(buttonID)
	}
	if cfg.ButtonsConfig.Back != nil {
		buttonID := m.customButtonID("back")
		cfg.DiscordConfig.RemoveComponentHandler(buttonID)
	}
	if cfg.ButtonsConfig.Stop != nil {
		buttonID := m.customButtonID("stop")
		cfg.DiscordConfig.RemoveComponentHandler(buttonID)
	}
	if cfg.ButtonsConfig.Next != nil {
		buttonID := m.customButtonID("next")
		cfg.DiscordConfig.RemoveComponentHandler(buttonID)
	}
	if cfg.ButtonsConfig.Last != nil {
		buttonID := m.customButtonID("last")
		cfg.DiscordConfig.RemoveComponentHandler(buttonID)
	}
	log.WithFields(log.Fields{"paginator": m.id}).Trace("deregistered component handlers")
}

// itemsPerPage returns the number of items per page. If the
// ItemsPerPage field is 0, it returns the default number of items
// per page.
func (m *message) getItemsPerPage() int {
	return m.paginator.config.ItemsPerPage
}

// hasExpired returns true if the paginator has expired.
func (m *message) hasExpired() bool {
	return !m.expiry.IsZero() && m.expiry.Before(time.Now())
}

// pageThroughItems is called when a page button is selected in a paginated message.
func pageThroughItems(s *discordgo.Session, i *discordgo.InteractionCreate) {
	ids := strings.Split(i.Interaction.MessageComponentData().CustomID, ":")
	paginatorID, messageID, action := ids[0], ids[1], ids[2]

	manager.mutex.Lock()
	paginator, ok := manager.paginators[paginatorID]
	manager.mutex.Unlock()
	if !ok {
		log.WithFields(log.Fields{"paginator": paginatorID}).Error("paginator not found")
		return
	}
	paginator.mutex.Lock()
	defer paginator.mutex.Unlock()
	m, ok := paginator.messages[messageID]
	if !ok {
		return
	}

	switch action {
	case "first":
		m.currentPage = 0

	case "back":
		m.currentPage--

	case "next":
		m.currentPage++

	case "last":
		m.currentPage = m.pageCount() - 1
	}

	m.expiry = time.Now().Add(m.paginator.config.IdleWait)
	m.editInteractionResponse(s, i)
}

// customButtonID returns the custom ID for a button in the paginator.
func (m *message) customButtonID(buttonText string) string {
	return fmt.Sprintf("%s:%s:%s", m.paginator.id, m.id, buttonText)
}
