package db

type User struct {
	ID                int
	Username          string
	ProfilePictureURL string
	Friends           []User
	Chats             []Chat
}
