package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
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
	Identifier    string          `json:"identifier"`
	Commands      []CommandObject `json:"commands"`
}

// Response details
type ServerResponse struct {
	Message string
	Details string
}

func startBot(token string, botName string, targetChannel string, commands []CommandObject) {
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

	//TODO: Check for spam/message frequency

	for {
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
			conn.Write([]byte("PONG " + msgParts[1]))
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

func main() {
	h := http.NewServeMux()

	h.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "A-MEI-zing!")
	})

	h.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "Raising my APM!")
	})

	h.HandleFunc("/createBot", func(w http.ResponseWriter, r *http.Request) {
		/* TODO: Check for ->
		// command format (not empty) (return error stating that the commands are empty)
		// check for existing bots for identifier/signature (limit one per user)
		// limit time for bot run time (2-6 hours)
		// allow a user to kill a bot they created
		*/
		var botInfo BotInfo
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}
		err := json.NewDecoder(r.Body).Decode(&botInfo)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		go startBot(botInfo.Token, botInfo.BotName, botInfo.TargetChannel, botInfo.Commands)
	})

	h.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		res := ServerResponse{Message: "Route error", Details: "No endpoint at this location"}
		json.NewEncoder(w).Encode(res)
	})

	err := http.ListenAndServe(":7000", h)
	log.Println("Listening...")
	log.Fatal(err)
}
