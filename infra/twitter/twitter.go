package twitter

import (
	"time"
)

type Tweet struct {
	ID        string    `json:"id_str"`
	User      *User     `json:"user"`
	Text      string    `json:"text"`
	CreatedAt time.Time `json:"create_at"`
}

type User struct {
	ID         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}
