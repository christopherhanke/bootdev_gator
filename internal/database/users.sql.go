// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.27.0
// source: users.sql

package database

import (
	"context"
	"database/sql"

	"github.com/google/uuid"
)

const createUser = `-- name: CreateUser :one
insert into users (id, created_at, updated_at, name)
values (
    $1,
    $2,
    $3,
    $4
)
returning id, created_at, updated_at, name
`

type CreateUserParams struct {
	ID        uuid.UUID
	CreatedAt sql.NullTime
	UpdatedAt sql.NullTime
	Name      string
}

func (q *Queries) CreateUser(ctx context.Context, arg CreateUserParams) (User, error) {
	row := q.db.QueryRowContext(ctx, createUser,
		arg.ID,
		arg.CreatedAt,
		arg.UpdatedAt,
		arg.Name,
	)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}

const getUsers = `-- name: GetUsers :one
select id, created_at, updated_at, name from users where name = $1 limit 1
`

func (q *Queries) GetUsers(ctx context.Context, name string) (User, error) {
	row := q.db.QueryRowContext(ctx, getUsers, name)
	var i User
	err := row.Scan(
		&i.ID,
		&i.CreatedAt,
		&i.UpdatedAt,
		&i.Name,
	)
	return i, err
}
