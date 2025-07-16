package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	apperr "github.com/seacite-tech/compendium/user-service/internal/error"
	"github.com/seacite-tech/compendium/user-service/internal/model"
	"github.com/ztrue/tracerr"
)

type PgUserRepository struct {
	db *sql.DB
}

func NewPgUserRepository(db *sql.DB) *PgUserRepository {
	return &PgUserRepository{
		db: db,
	}
}

func (r *PgUserRepository) FindByEmail(ctx context.Context, email string) (*model.User, error) {
	user := &model.User{}
	query := `
		SELECT id, name, email, is_email_verified, is_admin, password_hash, created_at
		FROM users
		WHERE email = $1
	`
	row := r.db.QueryRowContext(ctx, query, email)

	err := row.Scan(
		&user.Id,
		&user.Name,
		&user.Email,
		&user.IsEmailVerified,
		&user.IsAdmin,
		&user.PasswordHash,
		&user.CreatedAt,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Wrap(err)
	}

	return user, nil
}

func (r *PgUserRepository) UpdateIsEmailVerifiedByEmail(ctx context.Context, email string, isEmailVerified bool) error {
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

func (r *PgUserRepository) UpdatePasswordHash(ctx context.Context, id uuid.UUID, passwordHash []byte) error {
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

func (r *PgUserRepository) UpdatePasswordHashAndCreatedAt(ctx context.Context, id uuid.UUID, passwordHash []byte, createdAt time.Time) error {
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

func (r *PgUserRepository) CreateUser(ctx context.Context, user model.User) error {
	query := `
		INSERT INTO users (id, name, email, is_email_verified, is_admin, password_hash, created_at)
		VALUES ($1, $2, $3, $4, $5, $6, $7)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		user.Id,
		user.Name,
		user.Email,
		user.IsEmailVerified,
		user.IsAdmin,
		user.PasswordHash,
		user.CreatedAt,
	)

	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" { // 23505 is the SQLSTATE for unique_violation
			return apperr.Errorf(apperr.EmailTakenError, "email is already taken")
		}

		return tracerr.Wrap(err)
	}

	return nil
}
