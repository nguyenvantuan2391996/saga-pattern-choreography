package model

import "time"

type Order struct {
	ID        int32     `json:"id"`
	UserID    int32     `json:"user_id"`
	Status    string    `json:"status"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
