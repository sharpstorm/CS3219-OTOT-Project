package model

import "time"

type ApiToken struct {
	Id        int       `bun:"token_id"`
	Token     string    `bun:"token"`
	IsEnabled bool      `bun:"is_enabled"`
	CreatedAt time.Time `bun:"created_at"`
}
