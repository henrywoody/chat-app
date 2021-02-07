package main

import (
	"fmt"
	"sync"
	"time"
)

type Database struct {
	data                  map[string]*Room
	dataLock              sync.RWMutex
	lastMessageID         int
	messageIDLock         sync.Mutex
	roomSubscriptions     map[string]map[string]func([]*Message)
	roomSubscriptionsLock sync.RWMutex
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
	db := &Database{
		data:              make(map[string]*Room),
		roomSubscriptions: make(map[string]map[string]func([]*Message)),
	}

	db.CreateRoom("General")
	db.CreateRoom("Random")

	return db
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
	db.roomSubscriptions[name] = make(map[string]func([]*Message))
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

	db.roomSubscriptionsLock.RLock()
	defer db.roomSubscriptionsLock.RUnlock()
	subscriptions, _ := db.roomSubscriptions[roomName]
	for _, subscription := range subscriptions {
		subscription(room.Messages)
	}

	return message, nil
}

func (db *Database) SubscribeToRoom(roomName string, f func([]*Message)) string {
	db.roomSubscriptionsLock.Lock()
	defer db.roomSubscriptionsLock.Unlock()

	nAttempts := 0
	id := fmt.Sprintf("%d-%d", time.Now().UnixNano(), nAttempts)
	for _, isIDTaken := db.roomSubscriptions[roomName][id]; isIDTaken; nAttempts++ {
		id = fmt.Sprintf("%d-%d", time.Now().UnixNano(), nAttempts)
	}

	db.roomSubscriptions[roomName][id] = f

	return id
}

func (db *Database) UnsubscribeFromRoom(roomName string, id string) {
	delete(db.roomSubscriptions[roomName], id)
}
