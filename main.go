package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
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

type commandObject struct {
	User     string
	Command  string
	Response string
}

func main() {
	// Connect to the twitch server
	conn, err := net.Dial("tcp", "irc.chat.twitch.tv:6667")
	if err != nil {
		panic(err)
	}

	// Token, username, channel
	conn.Write([]byte("PASS " + "oauth:" + os.Getenv("BOT_OAUTH_TOKEN") + "\r\n"))
	conn.Write([]byte("NICK " + os.Getenv("BOT_ACCOUNT_NAME") + "\r\n"))
	conn.Write([]byte("JOIN " + os.Getenv("CHANNEL_NAME") + "\r\n"))
	defer conn.Close()

	// Handles reading from the connection
	tp := textproto.NewReader(bufio.NewReader(conn))

	//TODO: Check for spam/message frequency

	filePath := "./commands.json"
	fmt.Printf("// Reading file %s\n", filePath)
	file, fileReadError := ioutil.ReadFile(filePath)
	if fileReadError != nil {
		fmt.Printf("// Error while reading file %s\n", filePath)
		fmt.Printf("File error: %v\n", fileReadError)
		os.Exit(1)
	}

	var commands []commandObject

	readJSONError := json.Unmarshal(file, &commands)
	if readJSONError != nil {
		fmt.Println("Error:", readJSONError)
		os.Exit(1)
	}

	for {
		msg, err := tp.ReadLine()
		if err == io.EOF {
			continue
		} else if err != nil {
			panic(err)
		}

		// Split the msg by spaces
		msgParts := strings.Split(msg, " ")

		fmt.Println(msgParts)

		// Respond with PONG
		if msgParts[0] == "PING" {
			conn.Write([]byte("PONG " + msgParts[1]))
			continue
		}

		// If msg contains PRIVMSG then respond
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
