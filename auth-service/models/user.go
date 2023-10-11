package models

import (
	"database/sql"
	dbhelper "github.com/Noah-Wilderom/video-buffer/shared/helpers/db"
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
	id, err := dbhelper.Insert(c.DB)
}
