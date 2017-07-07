package main

import (
	"encoding/json"
	"log"
	"mightbot/botlogic"
	"net/http"
)

type ChannelMessage struct {
	Message string `json:"message"`
	UUID    string `json:"uuid"`
}

// Structure of bot creator's submission
type BotInfo struct {
	Token         string                   `json:"token"`
	BotName       string                   `json:"botName"`
	TargetChannel string                   `json:"targetChannel"`
	UUID          string                   `json:"uuid"`
	Commands      []botlogic.CommandObject `json:"commands"`
}

// Response details
type ServerResponse struct {
	Message string
	Details string
}

func decodeValidate(w http.ResponseWriter, r *http.Request, msgObj interface{}) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}

	err := json.NewDecoder(r.Body).Decode(&msgObj)
	if err != nil {
		http.Error(w, err.Error(), 400)
		return
	}
}

func main() {
	botDirectory := []botlogic.BotRecord{}
	serve := http.NewServeMux()

	serve.HandleFunc("/refresh", func(w http.ResponseWriter, r *http.Request) {
		var refreshInfo ChannelMessage
		decodeValidate(w, r, &refreshInfo)

		for _, v := range botDirectory {
			if v.UUID == refreshInfo.UUID {
				v.BotChannel <- "refresh"
			}
		}
	})

	serve.HandleFunc("/close", func(w http.ResponseWriter, r *http.Request) {
		var closeInfo ChannelMessage
		decodeValidate(w, r, &closeInfo)

		for _, v := range botDirectory {
			if v.UUID == closeInfo.UUID {
				v.BotChannel <- "close"
			}
		}
	})

	serve.HandleFunc("/createBot", func(w http.ResponseWriter, r *http.Request) {
		var botExists bool
		var botInfo BotInfo
		decodeValidate(w, r, &botInfo)

		for _, val := range botDirectory {
			if val.UUID == botInfo.UUID {
				botExists = true
			}
		}

		if botExists == true {
			return
		}

		newChannel := make(chan string)
		newRecord := botlogic.BotRecord{UUID: botInfo.UUID, BotChannel: newChannel}
		botDirectory = append(botDirectory, newRecord)
		go botlogic.StartBot(botInfo.Token, botInfo.BotName, botInfo.TargetChannel, botInfo.Commands, newRecord)
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
