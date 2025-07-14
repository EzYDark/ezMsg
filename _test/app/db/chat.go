package db

type Chat struct {
	ID       int
	Members  []*User
	Messages []*Message
	Files    []*File
}
