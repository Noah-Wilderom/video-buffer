package models

import "database/sql"

type (
	Models struct {
		User UserConn
	}
)

var db *sql.DB

func New(dbPool *sql.DB) Models {
	db = dbPool

	return Models{
		User: UserConn{
			DB: dbPool,
		},
	}
}
