package models

import "time"

type Country struct {
	ID        int       `db:"id"`
	Name      string    `db:"name"`
	Code      string    `db:"code"`
	Currency  string    `db:"currency"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}
