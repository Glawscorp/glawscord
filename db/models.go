package db

import (
	"time"
)

type User struct {
	ID       int       `json:"id" db:"id"`
	Username string    `json:"username" db:"username"`
	Password string    `json:"password" db:"password"`
	JoinedAt time.Time `json:"joined_at" db:"joined_at"`
}

type UserMessage struct {
	ID       int       `json:"id"`
	Sender   int       `json:"sender"`
	Receiver int       `json:"receiver"`
	SentAt   time.Time `json:"sent_at"`
	Content  string    `json:"content"`
}
