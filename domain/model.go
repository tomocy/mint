package domain

import (
	"fmt"
	"time"
)

type Tweet struct {
	ID        string
	User      *User
	Text      string
	CreatedAt time.Time
}

func (t *Tweet) String() string {
	return fmt.Sprintf("%s @%s\t%s\n\n\t%s\t\n\n", t.User.Name, t.User.ScreenName, t.CreatedAt, t.Text)
}

type User struct {
	ID         string
	Name       string
	ScreenName string
}
