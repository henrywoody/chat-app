package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"

	"golang.org/x/net/websocket"
)

func main() {
	db := NewDatabase()

	mux := http.NewServeMux()
	mux.Handle("/", HandleStatic("index.html", "text/html"))
	mux.Handle("/static/css/", HandleStatic("", "text/css"))
	mux.Handle("/static/js/", HandleStatic("", "application/javascript"))
	mux.HandleFunc("/rooms", HandleRooms)
	mux.HandleFunc("/rooms/", HandleGetRoom)
	mux.Handle("/messages", websocket.Handler(HandleSocket(db)))
	muxWithMiddleware := DBMiddleware(db)(mux)

	port := "8080"
	log.Printf("Server is running on port %s", port)
	http.ListenAndServe(":"+port, muxWithMiddleware)
}

type WebsocketConn struct {
	conn               *websocket.Conn
	roomName           string
	roomSubscriptionID string
}

func HandleSocket(db *Database) func(*websocket.Conn) {
	return func(ws *websocket.Conn) {
		c := &WebsocketConn{conn: ws}
		SetSocketRoom(c, db)
		go ReadSocket(c, db)
		go WriteSocket(c, db)
		for {
		}
	}
}

type MessageCreateInput struct {
	SenderName string `json:"senderName"`
	Body       string `json:"body"`
}

type RoomSelection struct {
	SelectRoom string `json:"selectRoom"`
}

func SetSocketRoom(c *WebsocketConn, db *Database) {
	for {
		var roomSelection RoomSelection
		if err := websocket.JSON.Receive(c.conn, &roomSelection); err != nil {
			log.Fatal("error encountered while reading room selection: ", err)
		}
		if roomSelection.SelectRoom != "" {
			c.roomName = roomSelection.SelectRoom
			return
		}
	}
}

func ReadSocket(c *WebsocketConn, db *Database) {
	for {
		var input MessageCreateInput
		if err := websocket.JSON.Receive(c.conn, &input); err != nil {
			if err.Error() == "EOF" {
				c.conn.Close()
				db.UnsubscribeFromRoom(c.roomName, c.roomSubscriptionID)
				return
			}
			log.Fatal("error encountered while reading: ", err)
		}
		dbInput := &Message{SenderName: input.SenderName, Body: input.Body}
		db.CreateMessage(c.roomName, dbInput)
	}
}

func WriteSocket(c *WebsocketConn, db *Database) {
	writeMessages := func(messages []*Message) {
		err := websocket.JSON.Send(c.conn, messages)
		if err != nil {
			log.Fatal("error encountered while writing: ", err)
		}
	}

	room, _ := db.GetRoom(c.roomName)
	writeMessages(room.Messages)

	c.roomSubscriptionID = db.SubscribeToRoom(c.roomName, writeMessages)
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
