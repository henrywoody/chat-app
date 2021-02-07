package main

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
	data          map[string]*Room
	dataLock      sync.RWMutex
	lastMessageID int
	messageIDLock sync.Mutex
}

type Room struct {
	Name     string     `json:"name"`
	Messages []*Message `json:"messages"`
}

type Message struct {
	ID         int       `json:"ID"`
	SenderName string    `json:"senderName"`
	SentAt     time.Time `json:"sentAt"`
	Body       string    `json:"body"`
}

func NewDatabase() *Database {
	return &Database{
		data: make(map[string]*Room),
	}
}

func (db *Database) GetRooms() []*Room {
	db.dataLock.RLock()
	defer db.dataLock.RUnlock()

	rooms := make([]*Room, 0, len(db.data))
	for _, room := range db.data {
		roomWithoutMessages := &Room{Name: room.Name}
		rooms = append(rooms, roomWithoutMessages)
	}
	return rooms
}

func (db *Database) GetRoom(name string) (*Room, error) {
	db.dataLock.RLock()
	defer db.dataLock.RUnlock()

	room, ok := db.data[name]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func (db *Database) CreateRoom(name string) (*Room, error) {
	db.dataLock.Lock()
	defer db.dataLock.Unlock()

	_, isNameTaken := db.data[name]
	if isNameTaken {
		return nil, fmt.Errorf("room name is taken")
	}
	room := &Room{Name: name, Messages: []*Message{}}
	db.data[name] = room
	return room, nil
}

func (db *Database) CreateMessage(roomName string, message *Message) (*Message, error) {
	db.dataLock.Lock()
	defer db.dataLock.Unlock()

	room, ok := db.data[roomName]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}

	db.messageIDLock.Lock()
	defer db.messageIDLock.Unlock()
	db.lastMessageID++

	message.ID = db.lastMessageID
	message.SentAt = time.Now()
	room.Messages = append(room.Messages, message)
	return message, nil
}
