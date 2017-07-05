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
	"time"
)

type BotRecord struct {
	BotChannel chan string
	UUID       string
}

type ChannelMessage struct {
	Message string `json:"message"`
	UUID    string `json:"uuid"`
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

// Response details
type ServerResponse struct {
	Message string
	Details string
}

func checkForCommand(message string) string {
	var commandMessage []string
	if strings.HasPrefix(message, "!") {
		commandMessage = strings.SplitAfter(message, "!")
	}

	return strings.Join(commandMessage, "")
}

func startBot(token string, botName string, targetChannel string, commands []CommandObject, botRecord BotRecord) {
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

	ticker := time.NewTicker(time.Minute)
	go func() {
		for t := range ticker.C {
			routineTimer++
			log.Println(t)
		}
	}()

	go func() {
		for message := range botRecord.BotChannel {
			switch message {
			case "refresh":
				routineTimer = 0
			case "close":
				return
			default:
			}
		}
	}()

	for {
		if routineTimer >= 30 {
			return
		}

		msg, err := tp.ReadLine()
		if err == io.EOF {
			fmt.Println("EOF", err)
			continue
		} else if err != nil {
			panic(err)
		}

		msgParts := strings.Split(msg, " ")

		// For logging purposes
		fmt.Println(msgParts)
		fmt.Println(msgParts[1])
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

func main() {
	botDirectory := []BotRecord{}
	serve := http.NewServeMux()

	serve.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		var refreshInfo ChannelMessage
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&refreshInfo)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		for _, v := range botDirectory {
			if v.UUID == refreshInfo.UUID {
				v.BotChannel <- "refresh"
			}
		}
	})

	serve.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var closeInfo ChannelMessage
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&closeInfo)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		for _, v := range botDirectory {
			if v.UUID == closeInfo.UUID {
				v.BotChannel <- "close"
			}
		}
	})

	serve.HandleFunc("/createBot", func(w http.ResponseWriter, r *http.Request) {
		/* TODO: Check for ->
		// command format (not empty) (return error stating that the commands are empty)
		// check for existing bots for UUID/signature (limit one per user)
		// limit for bot run time 30 minutes, refresh pings reset timer to 0
		// allow a user to kill a bot they created
		*/
		var botInfo BotInfo
		var botExists bool
		if r.Body == nil {
			http.Error(w, "Please send a request body", 400)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&botInfo)
		if err != nil {
			http.Error(w, err.Error(), 400)
			return
		}

		for _, val := range botDirectory {
			if val.UUID == botInfo.UUID {
				botExists = true
			}
		}

		if botExists == true {
			return
		}

		newChannel := make(chan string)
		newRecord := BotRecord{UUID: botInfo.UUID, BotChannel: newChannel}
		botDirectory = append(botDirectory, newRecord)
		go startBot(botInfo.Token, botInfo.BotName, botInfo.TargetChannel, botInfo.Commands, newRecord)
	})

	serve.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(404)
		res := ServerResponse{Message: "Route error", Details: "No endpoint at this location"}
		json.NewEncoder(w).Encode(res)
	})

	err := http.ListenAndServe(":7000", serve)
	log.Println("Listening...")
	log.Fatal(err)
}
