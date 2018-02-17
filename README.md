# MightBot
A simple distributed Twitch chatbot written in Go, which can be used to create and run multiple chatbots at once.

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
1. There's an amount of debounce built into the server, so the bot will respond to messages, but then wait for a time after that. This is sort of good, in that it prevents spam by default, but could cause issues, so I do plan on looking into where this is happening.

## TODO
1. Add ability to capture type of command (contains or exact), and user name of message origin.
This would be used to target users for bans when they use certain words or create games that track chat user progress through a task or group task.

2. Improve error handling, (ie, a user tries to refresh a closed/expired bot).

3. Do actual tests in the case of having multple bots running at the same time.
