package db

type MsgStatus uint8

const (
	Sent MsgStatus = iota
	Received
)
