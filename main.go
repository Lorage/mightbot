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

func remove(s []botlogic.BotRecord, i int) []botlogic.BotRecord {
	s[len(s)-1], s[i] = s[i], s[len(s)-1]
	return s[:len(s)-1]
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

		for index, bot := range botDirectory {
			if bot.UUID == closeInfo.UUID {
				bot.BotChannel <- "close"
				botDirectory = remove(botDirectory, index)
			}
		}
	})

	serve.HandleFunc("/createBot", func(w http.ResponseWriter, r *http.Request) {
		var botExists bool
		var botInfo botlogic.BotInfo
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
		go botlogic.StartBot(&botInfo, newRecord)
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
