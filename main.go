package main

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
	"os"
	"strings"
)

func checkForCommand(message string) string {
	var commandMessage []string
	if strings.HasPrefix(message, "!") {
		commandMessage = strings.SplitAfter(message, "!")
	}

	return strings.Join(commandMessage, "")
}

func removeIndex(s []string, index int) []string {
	return append(s[:index], s[index+1:]...)
}

func main() {
	// connect to the twitch server
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	// token, username, channel
	conn.Write([]byte("PASS " + "oauth:" + os.Getenv("BOT_OAUTH_TOKEN") + "\r\n"))
	conn.Write([]byte("NICK " + os.Getenv("BOT_ACCOUNT_NAME") + "\r\n"))
	conn.Write([]byte("JOIN " + os.Getenv("CHANNEL_NAME") + "\r\n"))
	defer conn.Close()

	// handles reading from the connection
	tp := textproto.NewReader(bufio.NewReader(conn))

	for {
		msg, err := tp.ReadLine()
		if err != nil {
			panic(err)
		}

		// split the msg by spaces
		msgParts := strings.Split(msg, " ")

		fmt.Println(msgParts)

		// if the msg contains PING you're required to
		// respond with PONG else you get kicked
		if msgParts[0] == "PING" {
			conn.Write([]byte("PONG " + msgParts[1]))
			continue
		}

		// if msg contains PRIVMSG then respond
		if msgParts[1] == "PRIVMSG" {
			// echo back the same message
			messageText := strings.Split(msg, os.Getenv("CHANNEL_NAME")+" :")
			var message = []byte("PRIVMSG " + os.Getenv("CHANNEL_NAME") + " :" + messageText[1] + "\r\n")
			conn.Write(message)
		}
	}
}
