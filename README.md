# DisGoPage

DisGoPage is a Discord pagination library for Go, built on top of [DiscordGo](https://github.com/bwmarrin/discordgo). It provides an easy way to create paginated messages in Discord, with customizable navigation buttons and configuration options.

## Features

- Create paginated messages with embeds
- Support for both regular messages and interaction responses
- Customizable navigation buttons (First, Back, Next, Last)
- Automatic cleanup of expired messages
- Configurable items per page
- Customizable embed colors
- Idle timeout configuration

## Installation

```bash
go get github.com/rbrabson/disgopage
```

## Usage

### Basic Example

```go
package main

import (
    "github.com/bwmarrin/discordgo"
    "github.com/rbrabson/disgopage"
)

func main() {
    // Create a new Discord session
    dg, err := discordgo.New("Bot " + "YOUR_BOT_TOKEN")
    if err != nil {
        // Handle error
    }
    
    // Create a new paginator
    p := disgopage.NewPaginator(
        disgopage.WithDiscordConfig(
            disgopage.DiscordConfig{
                Session:                dg,
                AddComponentHandler:    addComponentHandler,
                RemoveComponentHandler: removeComponentHandler,
            },
        ),
    )
    
    // Create embed fields for pagination
    embedFields := []*discordgo.MessageEmbedField{
        {
            Name:  "Field 1",
            Value: "Value 1",
        },
        {
            Name:  "Field 2",
            Value: "Value 2",
        },
        // Add more fields...
    }
    
    // Send a paginated message
    err = p.CreateMessage(dg, channelID, "Paginated Message Title", embedFields)
    if err != nil {
        // Handle error
    }
}

// Component handler functions
func addComponentHandler(key string, handler func(*discordgo.Session, *discordgo.InteractionCreate)) {
    // Add component handler to your bot
}

func removeComponentHandler(key string) {
    // Remove component handler from your bot
}
```

### Slash Command Example

```go
// Handle a slash command interaction
func paginatorCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
    p := disgopage.NewPaginator(
        disgopage.WithDiscordConfig(
            disgopage.DiscordConfig{
                Session:                s,
                AddComponentHandler:    addComponentHandler,
                RemoveComponentHandler: removeComponentHandler,
            },
        ),
    )
    
    // Create embed fields
    embedFields := []*discordgo.MessageEmbedField{
        // Your embed fields...
    }
    
    // Create an ephemeral paginated response
    err := p.CreateInteractionResponse(s, i, "Paginated Response", embedFields, true)
    if err != nil {
        // Handle error
    }
}
```

## Configuration Options

DisGoPage provides several configuration options:

```go
// Customize the paginator
p := disgopage.NewPaginator(
    // Set custom button styles and emojis
    disgopage.WithButtonsConfig(disgopage.ButtonsConfig{
        First: &disgopage.ComponentOption{
            Emoji: &discordgo.ComponentEmoji{Name: "⏮️"},
            Style: discordgo.PrimaryButton,
        },
        // Configure other buttons...
    }),
    
    // Set custom embed color
    disgopage.WithEmbedColor(0x4c50c1),
    
    // Set items per page
    disgopage.WithItemsPerPage(5),
    
    // Set idle timeout
    disgopage.WithIdleWait(time.Minute * 5),
    
    // Set custom ID prefix
    disgopage.WithCustomIDPrefix("my-paginator"),
)
```

## Complete Example

See the [examples](https://github.com/rbrabson/disgopage/tree/main/examples) directory for a complete working example.

## License

This project is licensed under the GNU General Public License v3.0 - see the [LICENSE](LICENSE) file for details.