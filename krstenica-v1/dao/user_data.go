package dao

import "time"

type HramDo struct {
	HramID    uint      `json:"hram_id"`
	HramName  string    `json:"hram_name"`
	CreatedAt time.Time `json:"created_at"`
}
