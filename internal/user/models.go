package user

import "github.com/uptrace/bun"

type User struct {
	bun.BaseModel `bun:"table:users,alias:u"`

	Id     int64  `bun:",pk,autoincrement"`
	Secret string `bun:", not null,unique"`
}
