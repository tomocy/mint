package domain

import "time"

type Tweet struct {
	ID        string
	Text      string
	CreatedAt time.Time
}
