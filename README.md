# MightBot
A simple Twitch chatbot written in Go.

## Setup
1. Run `go build` & `./mightbot`
2. Send the token, channel name, uuid, commands, and botname to `/createBot`
3. Commands are in the form of:
``` 
{ "command": [what to listen for], "response": "[your bot's response]" } 
```

### Closing channels
Bots close automatically after 30 minutes, unless you send the `refresh` message to `/refresh`.

Bots can be closed specifically by sending a matching `uuid` to a created bot and the `close` message to `/close`.

### OAUTH
Use [Twitchapps.com](http://twitchapps.com/tmi/) to generate the OAUTH token.

## Issues
1. Channel logic doesn't work as expected. Messages sent to the channel only trigger in the for loop after a message is recieved on the net/http `dial` Reader.

## TODO
1. Rewrite the method/signature for refreshing and closing bot goroutines.
2. Improve error handling.
