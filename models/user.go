package models

type User struct {
	ID      int64  `json:"id"`
	Name    string `json:"name"`
	Balance string `json:"balance"`
}
