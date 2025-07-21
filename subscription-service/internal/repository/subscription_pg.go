package repository

import (
	"database/sql"
	"fmt"

	"github.com/adslmgrv/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type pgSubscriptionRepository struct {
	db *sql.DB
}

func NewPgSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &pgSubscriptionRepository{db: db}
}

func (r *pgSubscriptionRepository) CreateSubscription(sub model.Subscription) error {
	query := `
		INSERT INTO subscriptions (user_id, subscription_level, till, since)
		VALUES ($1, $2, $3, $4)`

	err := r.db.QueryRow(query, sub.UserID, sub.SubscriptionLevel, sub.Till, sub.Since).Scan(&sub.UserID)
	if err != nil {
		return tracerr.Errorf("failed to create subscription: %w", err)
	}

	return nil
}

func (r *pgSubscriptionRepository) GetSubscriptionByUserID(userID uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT user_id, subscription_level, till, since
		FROM subscriptions
		WHERE user_id = $1`

	sub := &model.Subscription{}
	err := r.db.QueryRow(query, userID).Scan(
		&sub.UserID,
		&sub.SubscriptionLevel,
		&sub.Till,
		&sub.Since,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}

		return nil, tracerr.Errorf("failed to get subscription by user ID %s: %w", userID, err)
	}
	return sub, nil
}

func (r *pgSubscriptionRepository) RemoveSubscription(userID uuid.UUID) error {
	query := `
		DELETE FROM subscriptions
		WHERE user_id = $1`

	result, err := r.db.Exec(query, userID)
	if err != nil {
		return fmt.Errorf("failed to delete subscription for user ID %s: %w", userID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected after delete: %w", err)
	}

	if rowsAffected == 0 {
		return fmt.Errorf("no subscription found to delete for user ID %s", userID)
	}

	return nil
}
