# Bot -- Discord Bot client in Golang.
Bot program that streams messages from discord and lets you send messages to specific channels.</p>
Uses Discordgo package and bolt DB.</p>
**Version 1.0.0**
## Usage
Assuming you have Go installed and Go Bin path set, run:

```
go install
```

Usage:
  bot [command]

### Available Commands:
  - configure   Configure your bot
  - help        Help about any command
  - start       Start the bot to listen to messages

- Flags:</p>
  -h, --help   help for bot

Use "bot [command] --help" for more information about a command.

## After Start 
```
/addc  <channel name> to add new channel 
/listc to list channels
/default to set default channel
/send <channel name> <message> to send to specific channel
type in messages to send to default channel
```