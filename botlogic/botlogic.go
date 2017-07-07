package botlogic

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"net/textproto"
	"os"
	"strings"
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

func StartBot(token string, botName string, targetChannel string, commands []CommandObject, botRecord BotRecord) {
	var routineTimer int
	// Connect to the twitch server
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	// Token, username, channel
	conn.Write([]byte("PASS " + "oauth:" + token + "\r\n"))
	conn.Write([]byte("NICK " + botName + "\r\n"))
	conn.Write([]byte("JOIN " + targetChannel + "\r\n"))
	defer conn.Close()

	// Handles reading from the connection
	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		select {
		case channelResult := <-botRecord.BotChannel:
			switch {
			case channelResult == "refresh":
				routineTimer = 0
			case channelResult == "close":
				return
			}
		default:
			break
		}

		if routineTimer >= 30 {
			return
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
			messageText := strings.Split(msg, os.Getenv("CHANNEL_NAME")+" :")

			for command := range commands {
				var newString = strings.Join(messageText, "")
				if strings.Contains(newString, commands[command].Command) {
					newMessage = commands[command].Response
				}
			}

			var message = []byte("PRIVMSG " + os.Getenv("CHANNEL_NAME") + " :" + newMessage + "\r\n")
			conn.Write(message)
		}
	}
}
