package main

import "go.mongodb.org/mongo-driver/bson/primitive"


type Keyword struct {
	ID     primitive.ObjectID `bson:"_id,omitempty"`
	Keyword string            `bson:"keyword"`
	Reply   string            `bson:"reply"`
}

type IncomingMessage struct {
	Number  string `json:"number"`  // Nomor pengirim
	Message string `json:"message"` // Isi pesan
}