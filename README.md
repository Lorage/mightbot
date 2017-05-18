# MightBot
A simple Twitch chatbot written in Go.

#### Requires these environment variables
```bash
export CHANNEL_NAME="[your channel name]"
export BOT_ACCOUNT_NAME="[your bot account name]"
export BOT_OAUTH_TOKEN="[your bot OAuth token]"
```
## Setup
1. Create your environment variables in ~./bashrc
2. cd into the folder where main.go and commands.json live
3. Add new commands to commands.json in the form of:
``` 
{ "command": [what to listen for], "response": "[your bot's response]" } 
```
4. run `go run main.go`


#### OAUTH
Use [Twitchapps.com](http://twitchapps.com/tmi/) to generate the OAUTH token.

## TODO
1. Make sure there is no channel spam from the bot
2. Move from a JSON file based system to a simple GUI
3. Fix EOF error
