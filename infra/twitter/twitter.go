package twitter

type Tweet struct {
	ID        string `json:"id_str"`
	Text      string `json:"text"`
	CreatedAt string `json:"create_at"`
}

type User struct {
	ID         string `json:"id_str"`
	Name       string `json:"name"`
	ScreenName string `json:"screen_name"`
}
