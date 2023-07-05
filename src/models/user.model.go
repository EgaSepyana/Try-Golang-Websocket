package models

import "time"

type User struct {
	ID      string
	Addr    string
	EnterAt time.Time
}
