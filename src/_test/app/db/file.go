package db

import "time"

type File struct {
	ID         int
	Sender     *User
	URL        string
	UploadTime time.Time
}
