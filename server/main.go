package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
)

func main() {
	db := NewDatabase()

	mux := http.NewServeMux()
	mux.Handle("/", HandleStatic("index.html", "text/html"))
	mux.Handle("/static/css/", HandleStatic("", "text/css"))
	mux.Handle("/static/js/", HandleStatic("", "application/javascript"))
	mux.HandleFunc("/rooms", HandleRooms)
	mux.HandleFunc("/rooms/", HandleGetRoom)
	mux.HandleFunc("/messages", HandlePostMessage)
	muxWithMiddleware := DBMiddleware(db)(mux)

	port := "8080"
	log.Printf("Server is running on port %s", port)
	http.ListenAndServe(":"+port, muxWithMiddleware)
}

func HandleRooms(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		HandleGetRooms(w, r)
	case http.MethodPost:
		HandlePostRoom(w, r)
	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

func HandleGetRooms(w http.ResponseWriter, r *http.Request) {
	db := DBFromRequest(r)
	rooms := db.GetRooms()

	writeJSON(w, rooms)
}

func HandleGetRoom(w http.ResponseWriter, r *http.Request) {
	db := DBFromRequest(r)
	re := regexp.MustCompile(`^/rooms/(.+)$`)
	match := re.FindStringSubmatch(r.URL.Path)
	if match == nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Must supply a room name."))
		return
	}
	roomName := match[1]
	room, err := db.GetRoom(roomName)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}

	writeJSON(w, room)
}

type RoomCreateInput struct {
	Name string `json:"name"`
}

func HandlePostRoom(w http.ResponseWriter, r *http.Request) {
	var input RoomCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid format, must be JSON."))
		return
	}

	if input.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Key 'name' is required."))
		return
	}

	db := DBFromRequest(r)
	room, err := db.CreateRoom(input.Name)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	writeJSON(w, room)
}

type MessageCreateInput struct {
	RoomName   string `json:"roomName"`
	SenderName string `json:"senderName"`
	Body       string `json:"body"`
}

func HandlePostMessage(w http.ResponseWriter, r *http.Request) {
	var input MessageCreateInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid format, must be JSON."))
		return
	}

	if input.RoomName == "" || input.SenderName == "" || input.Body == "" {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Keys 'roomName', 'senderName', and 'body' are required."))
		return
	}

	dbInput := &Message{SenderName: input.SenderName, Body: input.Body}
	db := DBFromRequest(r)
	message, err := db.CreateMessage(input.RoomName, dbInput)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte(err.Error()))
		return
	}

	writeJSON(w, message)
}
