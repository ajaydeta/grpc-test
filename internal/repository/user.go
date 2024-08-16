package repository

import (
	"context"
	"database/sql"
	"errors"
	"tablelink/domain"
)

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

func (i *UserRepository) FindUserByEmail(ctx context.Context, email string) (*domain.User, error) {

	var user domain.User
	query := `
	SELECT 
	    id,
	    role_id,
	    email,
	    password,
	    name,
	    last_access
	FROM users WHERE email = $1
`

	err := i.db.
		QueryRowContext(ctx, query, email).
		Scan(
			&user.ID,
			&user.RoleID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.LastAccess,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (i *UserRepository) FindUserById(ctx context.Context, id int64) (*domain.User, error) {

	var user domain.User
	query := `
	SELECT 
	    id,
	    role_id,
	    email,
	    password,
	    name,
	    last_access
	FROM users WHERE id = $1
`

	err := i.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&user.ID,
			&user.RoleID,
			&user.Email,
			&user.Password,
			&user.Name,
			&user.LastAccess,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &user, nil
}

func (i *UserRepository) CreateUser(ctx context.Context, user *domain.User) error {
	query := `
	INSERT INTO users (role_id, email, password, name)
	VALUES ($1, $2, $3, $4)
`

	_, err := i.db.ExecContext(ctx, query, user.RoleID, user.Email, user.Password, user.Name)
	if err != nil {
		return err
	}

	return nil
}

func (i *UserRepository) UpdateNameUser(ctx context.Context, user *domain.User) error {
	query := `
	UPDATE users SET name = $1
	WHERE id = $2
`

	_, err := i.db.ExecContext(ctx, query, user.Name, user.ID)
	if err != nil {
		return err
	}

	return nil
}

func (i *UserRepository) FindRoleRightByRoleId(ctx context.Context, id int64) (*domain.RoleRight, error) {
	var role domain.RoleRight
	query := `
	SELECT 
	    id,
	    role_id,
	    r_create,
		r_read,
		r_update,
		r_delete
	FROM role_rights WHERE role_id = $1
`

	err := i.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&role.ID,
			role.RoleID,
			role.RCreate,
			role.RRead,
			role.RUpdate,
			role.RDelete,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (i *UserRepository) FindRoleById(ctx context.Context, id int64) (*domain.Role, error) {
	var role domain.Role
	query := `
	SELECT 
	    id,
	    name
	FROM roles WHERE id = $1
`

	err := i.db.
		QueryRowContext(ctx, query, id).
		Scan(
			&role.ID,
			&role.Name,
		)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, err
	}

	return &role, nil
}

func (i *UserRepository) DeleteUserById(ctx context.Context, id int64) error {
	query := `
	UPDATE users SET deleted_at = NOW()
	WHERE id = $1
`

	_, err := i.db.ExecContext(ctx, query, id)
	if err != nil {
		return err
	}

	return nil
}
