package main

import (
	"fmt"
	"time"
)

type Database struct {
	data map[string]*Room
}

type Room struct {
	Name     string     `json:"name"`
	Messages []*Message `json:"messages"`
}

type Message struct {
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
	rooms := make([]*Room, 0, len(db.data))
	for _, room := range db.data {
		roomWithoutMessages := &Room{Name: room.Name}
		rooms = append(rooms, roomWithoutMessages)
	}
	return rooms
}

func (db *Database) GetRoom(name string) (*Room, error) {
	room, ok := db.data[name]
	if !ok {
		return nil, fmt.Errorf("room not found")
	}
	return room, nil
}

func (db *Database) CreateRoom(name string) (*Room, error) {
	_, isNameTaken := db.data[name]
	if isNameTaken {
		return nil, fmt.Errorf("room name is taken")
	}
	room := &Room{Name: name, Messages: []*Message{}}
	db.data[name] = room
	return room, nil
}
