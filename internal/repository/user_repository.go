package repository

import (
	"context"
	"errors"

	apperror "github.com/hafiztri123/kki-be/internal/app_error"
	"github.com/hafiztri123/kki-be/internal/constants"
	"github.com/hafiztri123/kki-be/internal/models"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgxpool"
	"database/sql"
)

type UserRepository struct {
	db *pgxpool.Pool
}

func NewUserRepository(db *pgxpool.Pool) *UserRepository {
	return &UserRepository{
		db: db,
	}
}

func (r *UserRepository) InsertUser(ctx context.Context, user *models.User) error {

	sql := `INSERT INTO users (id, username, email, password, role, name, created_at, updated_at, deleted_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err := r.db.Exec(ctx, sql, 
		user.ID, 
		user.Username, 
		user.Email, 
		user.Password, 
		user.Role, 
		user.Name, 
		user.CreatedAt, 
		user.UpdatedAt, 
		user.DeletedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError

		if errors.As(err, &pgErr){
			if pgErr.Code == constants.UniqueConstraintViolationErrorCode {
				return apperror.ErrEmailAlreadyExists
			}
		}

		return err
	}

	return nil
}

func (r *UserRepository) GetUserByEmail(ctx context.Context, email string) (*models.User, error) {

	query := `SELECT id, username, email, password, role, name, created_at, updated_at, deleted_at FROM users WHERE email = $1`

	var user models.User
	err := r.db.QueryRow(ctx, query, email).Scan(
		&user.ID, 
		&user.Username, 
		&user.Email, 
		&user.Password, 
		&user.Role, 
		&user.Name, 
		&user.CreatedAt, 
		&user.UpdatedAt, 
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows)  {
			return nil, apperror.ErrNotFound
		}

		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) GetUsersByRole(ctx context.Context, role string, limit, offset int) ([]models.User, int64, error) {
	var totalCount int64
	countQuery := `SELECT COUNT(*) FROM users WHERE role = $1 AND deleted_at IS NULL`
	err := r.db.QueryRow(ctx, countQuery, role).Scan(&totalCount)
	if err != nil {
		return nil, 0, err
	}

	query := `SELECT id, username, email, password, role, name, created_at, updated_at, deleted_at
			  FROM users
			  WHERE role = $1 AND deleted_at IS NULL
			  ORDER BY created_at DESC
			  LIMIT $2 OFFSET $3`

	rows, err := r.db.Query(ctx, query, role, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var users []models.User
	for rows.Next() {
		var user models.User
		err := rows.Scan(
			&user.ID,
			&user.Username,
			&user.Email,
			&user.Password,
			&user.Role,
			&user.Name,
			&user.CreatedAt,
			&user.UpdatedAt,
			&user.DeletedAt,
		)
		if err != nil {
			return nil, 0, err
		}
		users = append(users, user)
	}

	return users, totalCount, nil
}

func (r *UserRepository) GetUserByID(ctx context.Context, id string) (*models.User, error) {
	query := `SELECT id, username, email, password, role, name, created_at, updated_at, deleted_at
			  FROM users
			  WHERE id = $1 AND deleted_at IS NULL`

	var user models.User
	err := r.db.QueryRow(ctx, query, id).Scan(
		&user.ID,
		&user.Username,
		&user.Email,
		&user.Password,
		&user.Role,
		&user.Name,
		&user.CreatedAt,
		&user.UpdatedAt,
		&user.DeletedAt,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, apperror.ErrNotFound
		}
		return nil, err
	}

	return &user, nil
}

func (r *UserRepository) UpdateUser(ctx context.Context, user *models.User) error {
	query := `UPDATE users
			  SET username = $1, email = $2, name = $3, updated_at = $4
			  WHERE id = $5 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query,
		user.Username,
		user.Email,
		user.Name,
		user.UpdatedAt,
		user.ID,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) {
			if pgErr.Code == constants.UniqueConstraintViolationErrorCode {
				return apperror.ErrEmailAlreadyExists
			}
		}
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

func (r *UserRepository) DeleteUser(ctx context.Context, id string) error {
	query := `UPDATE users
			  SET deleted_at = NOW()
			  WHERE id = $1 AND deleted_at IS NULL`

	result, err := r.db.Exec(ctx, query, id)
	if err != nil {
		return err
	}

	rowsAffected := result.RowsAffected()
	if rowsAffected == 0 {
		return apperror.ErrNotFound
	}

	return nil
}

