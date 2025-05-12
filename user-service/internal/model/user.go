package model

type User struct {
	ID           int64  `db:"id"`
	Name         string `db:"name"`
	Password     string `db:"password"`
	RefreshToken string `db:"refresh_token"`
	Role         string `db:"role"`
}
