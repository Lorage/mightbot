package botlogic

import (
	"bufio"
	"fmt"
	"net"
	"strings"
	"time"
)

type BotRecord struct {
	BotChannel chan string
	UUID       string
}

// Structure of []Commands
type CommandObject struct {
	Command  string
	Response string
}

// Structure of bot creator's submission
type BotInfo struct {
	Token         string          `json:"token"`
	BotName       string          `json:"botName"`
	TargetChannel string          `json:"targetChannel"`
	UUID          string          `json:"uuid"`
	Commands      []CommandObject `json:"commands"`
}

func StartBot(botDirectory *[]BotRecord, botInfo *BotInfo, botRecord BotRecord) {
	var lastPing = time.Now().Unix()
	// Connect to the twitch server
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	// Token, username, channel
	conn.Write([]byte("PASS " + "oauth:" + botInfo.Token + "\r\n"))
	conn.Write([]byte("NICK " + botInfo.BotName + "\r\n"))
	conn.Write([]byte("JOIN " + botInfo.TargetChannel + "\r\n"))
	defer conn.Close()

	go func() {
		scanner := bufio.NewScanner(bufio.NewReader(conn))
		for scanner.Scan() {
			msg := scanner.Text()
			msgParts := strings.Split(msg, " ")

			// For logging/debug purposes
			fmt.Println(msgParts)

			// Respond with PONG required
			if msgParts[0] == "PING" {
				conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
				continue
			}

			// Respond to channel member messages
			// Refactor into function and add new functionality
			if msgParts[1] == "PRIVMSG" {
				var newMessage string
				messageText := strings.Split(msg, botInfo.TargetChannel+" :")

				for command := range botInfo.Commands {
					var newString = strings.Join(messageText, "")
					if strings.Contains(newString, botInfo.Commands[command].Command) {
						newMessage = botInfo.Commands[command].Response
					}
				}

				var message = []byte("PRIVMSG " + botInfo.TargetChannel + " :" + newMessage + "\r\n")
				conn.Write(message)
			}
		}
	}()

	// Handles reading from the connection
	for {
		select {
		case channelResult := <-botRecord.BotChannel:
			switch {
			case channelResult == "refresh":
				fmt.Println("refreshed")
				lastPing = time.Now().Unix()
			case channelResult == "close":
				fmt.Println("closed")
				return
			}
		}

		// Checks if bot has been refreshed within last 30 minutes
		if time.Now().Unix()%lastPing > 1800 {
			for index, bot := range *botDirectory {
				if bot.UUID == botInfo.UUID {
					(*botDirectory)[index] = BotRecord{}
					return
				}
			}
		}
	}
}
