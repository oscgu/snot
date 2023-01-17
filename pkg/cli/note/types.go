package note

import "time"

type Note struct {
	Created     time.Time
	LastChanged time.Time
	Topic       string
	Author      string
	Title       string
	Content     string
}
