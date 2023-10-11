package models

import (
	"database/sql"
	"time"
)

type (
	User struct {
		Id        string    `json:"id"`
		Email     string    `json:"email"`
		Password  string    `json:"password"`
		CreatedAt time.Time `json:"created_at"`
		UpdatedAt time.Time `json:"updated_at"`
	}

	UserConn struct {
		DB *sql.DB
	}
)

func (c *UserConn) Create() {
	//
}
