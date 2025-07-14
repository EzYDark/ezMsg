package db

import (
	"time"

	"github.com/ezydark/ezlog/log"
)

type DB struct {
	RegisteredUsers []*User
}

var main_db *DB

func InitDB() *DB {
	if main_db != nil {
		log.Fatal().Msg("Database is already initialized")
	}

	var ezy_user User
	var kheper_user User

	ezy_user = User{
		ID:                123844,
		Username:          "EzY",
		ProfilePictureURL: "https://i.pinimg.com/736x/75/bc/fa/75bcfab6c24ba2cfeb6c16c1482c6b5f.jpg",
		Chats:             []*Chat{},
		Friends:           []*User{&kheper_user},
	}
	kheper_user = User{
		ID:                475774,
		Username:          "Kheper",
		ProfilePictureURL: "https://i.pinimg.com/736x/1e/43/a0/1e43a01262868b9b3a59912cb2e746f0.jpg",
		Chats:             []*Chat{},
		Friends:           []*User{&ezy_user},
	}

	file1 := File{
		ID:         1,
		Sender:     &kheper_user,
		URL:        "https://i.pinimg.com/originals/f0/bf/92/f0bf92400fdfb762aa109e1dfcfb4a9a.png",
		UploadTime: time.Now().Add(10 * time.Minute),
	}

	chat := Chat{
		ID:      778784,
		Members: []*User{&ezy_user, &kheper_user},
		Messages: []*Message{
			&Message{
				ID:        1,
				Sender:    &ezy_user,
				Message:   "Hello",
				Timestamp: time.Now().Add(time.Minute * 50),
				Status:    Received,
			},
			&Message{
				ID:        2,
				Sender:    &kheper_user,
				Message:   "Hi",
				Timestamp: time.Now().Add(time.Minute * 23),
				Status:    Received,
			},
			&Message{
				ID:        3,
				Sender:    &ezy_user,
				Message:   "How are you?",
				Timestamp: time.Now().Add(time.Minute * 4),
				Status:    Sent,
			},
			&Message{
				ID:        4,
				Sender:    &kheper_user,
				Message:   "I'm good, thanks!",
				Timestamp: time.Now().Add(time.Minute * 2),
				Status:    Received,
			},
			&Message{
				ID:        5,
				Sender:    &ezy_user,
				Message:   "That's good to hear!",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Sent,
			},
			&Message{
				ID:        6,
				Sender:    &kheper_user,
				Message:   "Nice to hear that!",
				Timestamp: time.Now(),
				Status:    Sent,
			},
			&Message{
				ID:        7,
				Sender:    &ezy_user,
				Message:   "What about you?",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Sent,
			},
			&Message{
				ID:        8,
				Sender:    &kheper_user,
				Message:   "I'm good too!",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Received,
			},
			&Message{
				ID:        9,
				Sender:    &ezy_user,
				Message:   "What you doin?",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Sent,
			},
			&Message{
				ID:        10,
				Sender:    &kheper_user,
				Message:   "I'm coding!",
				Files:     []*File{&file1},
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Received,
			},
			&Message{
				ID:        11,
				Sender:    &ezy_user,
				Message:   "Cool!",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Sent,
			},
			&Message{
				ID:        12,
				Sender:    &kheper_user,
				Message:   "I'm coding too!",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Received,
			},
			&Message{
				ID:        13,
				Sender:    &ezy_user,
				Message:   "Niceeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeeee!",
				Timestamp: time.Now().Add(time.Minute * 1),
				Status:    Sent,
			},
		},
	}

	ezy_user.Chats = append(ezy_user.Chats, &chat)
	kheper_user.Chats = append(kheper_user.Chats, &chat)

	main_db = &DB{
		RegisteredUsers: []*User{&ezy_user, &kheper_user},
	}

	return main_db
}

func GetDB() *DB {
	if main_db == nil {
		log.Fatal().Msg("Database is not initialized")
	}

	return main_db
}
