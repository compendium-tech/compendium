package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/compendium-tech/compendium/user-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ztrue/tracerr"
)

type pgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) UserRepository {
	return &pgUserRepository{
		db: db,
	}
}

func (r *pgUserRepository) GetUser(ctx context.Context, id uuid.UUID) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, name, email, is_admin, password_hash, created_at
		FROM users
		WHERE id = $1 AND is_email_verified = true
	`
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsAdmin,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, tracerr.Wrap(err)
	}

	user.IsEmailVerified = true

	return user, nil
}

func (r *pgUserRepository) FindUserByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, name, email, is_email_verified, is_admin, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsEmailVerified,
		&user.IsAdmin,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, tracerr.Wrap(err)
	}

	return user, nil
}

func (r *pgUserRepository) FindUserByVerifiedEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, name, email, is_admin, password_hash, created_at
		FROM users
		WHERE email = $1, is_email_verified = true
	`
	row := r.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsAdmin,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, tracerr.Wrap(err)
	}

	user.IsEmailVerified = true

	return user, nil
}

func (r *pgUserRepository) UpdateIsEmailVerifiedByEmail(ctx context.Context, email string, isEmailVerified bool) error {
	query := `
		UPDATE users
		SET is_email_verified = $1
		WHERE email = $2
	`
	res, err := r.db.ExecContext(ctx, query, isEmailVerified, email)
	if err != nil {
		return fmt.Errorf("failed to update is_email_verified for user %s: %w", email, err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no user found with email %s to update IsEmailVerified", email)
	}

	return nil
}

func (r *pgUserRepository) UpdateUserName(ctx context.Context, id uuid.UUID, name string) (*model.User, error) {
	query := `
		UPDATE users
		SET name = $1
		WHERE id = $2
		RETURNING id, name, email, is_email_verified, is_admin, password_hash, created_at
	`
	user := &model.User{}
	row := r.db.QueryRowContext(ctx, query, name, id)

	err := row.Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.IsEmailVerified,
		&user.IsAdmin,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, tracerr.Errorf("no user found with ID %s to update name", id)
		}

		return nil, tracerr.Wrap(err)
	}
	return user, nil
}

func (r *pgUserRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) error {
	query := `
		UPDATE users
		SET password_hash = $1
		WHERE email = $2
	`
	res, err := r.db.ExecContext(ctx, query, passwordHash, id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no user %s found to update PasswordHash and CreatedAt", id)
	}

	return nil
}

func (r *pgUserRepository) UpdatePasswordHashAndCreatedAt(ctx context.Context, id uuid.UUID, passwordHash []byte, createdAt time.Time) error {
	query := `
		UPDATE users
		SET password_hash = $1, created_at = $2
		WHERE id = $3
	`
	res, err := r.db.ExecContext(ctx, query, passwordHash, createdAt, id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no user %s found to update PasswordHash and CreatedAt", id)
	}

	return nil
}

func (r *pgUserRepository) CreateUser(ctx context.Context, user model.User, isEmailTaken *bool) error {
	query := `
		INSERT INTO users (id, name, email, is_email_verified, is_admin, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.ID,
		user.Name,
		user.Email,
		user.IsEmailVerified,
		user.IsAdmin,
		user.PasswordHash,
		user.CreatedAt,
	)

	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" { // 23505 is the SQLSTATE for unique_violation
			if isEmailTaken != nil {
				*isEmailTaken = true
			}
		}

		return tracerr.Wrap(err)
	}

	return nil
}
