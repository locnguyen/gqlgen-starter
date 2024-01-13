package models

import "time"

type Session struct {
	Token  string    `json:"token"`
	Expiry time.Time `json:"expiry"`
}
