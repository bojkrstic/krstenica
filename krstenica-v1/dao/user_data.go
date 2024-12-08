package dao

import "time"

type HramDo struct {
	HramID    uint      `json:"hram_id"`
	HramName  string    `json:"hram_name"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
}
