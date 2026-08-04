package models

import "time"

type Link struct {
	Short     string    `json:"Short"`
	Original  string    `json:"original"`
	UserID    string    `json:"user_id"`
	CratedAt  time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Clicks    int64     `json:"clicks"`
}
