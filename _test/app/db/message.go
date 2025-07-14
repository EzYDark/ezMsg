package db

import "time"

type Message struct {
	ID        int
	Sender    *User
	Message   string
	Files     []*File
	Timestamp time.Time
	Status    MsgStatus
}
