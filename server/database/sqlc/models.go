// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.24.0

package database

import (
	"github.com/jackc/pgx/v5/pgtype"
)

type Task struct {
	ID          pgtype.UUID
	Title       string
	Description pgtype.Text
}
