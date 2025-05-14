package db

import "time"

type Message struct {
	ID        int
	User      User
	Message   string
	Timestamp time.Time
	Status    MsgStatus
}
