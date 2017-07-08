package botlogic

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
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

	// Handles reading from the connection
	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		select {
		case channelResult := <-botRecord.BotChannel:
			switch {
			case channelResult == "refresh":
				lastPing = time.Now().Unix()
			case channelResult == "close":
				return
			}
		default:
			break
		}

		// TODO: time check doesn't clear botDirectory of bot botRecord 1800 is 30 minutes
		if time.Now().Unix()%lastPing > 1800 {
			//var blank = BotInfo{}
			for index, bot := range *botDirectory {
				if bot.UUID == botInfo.UUID {
					(*botDirectory)[index] = BotRecord{}
					return
				}
			}

		}

		msg, err := tp.ReadLine()
		if err == io.EOF {
			continue
		} else if err != nil {
			panic(err)
		}

		msgParts := strings.Split(msg, " ")

		// For logging purposes
		fmt.Println(msgParts)

		// Respond with PONG required
		if msgParts[0] == "PING" {
			conn.Write([]byte("PONG :tmi.twitch.tv\r\n"))
			continue
		}

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
}
