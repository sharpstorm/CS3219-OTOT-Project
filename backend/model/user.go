package model

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id       int    `bun:"user_id,pk,autoincrement"`
	Username string `bun:"user_name"`
	Password string `bun:"user_password"`
}
