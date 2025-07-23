package repository

import (
	"database/sql"

	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
	"github.com/ztrue/tracerr"
)

type pgSubscriptionRepository struct {
	db *sql.DB
}

func NewPgSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &pgSubscriptionRepository{db: db}
}

func (r *pgSubscriptionRepository) PutSubscription(sub model.Subscription) error {
	tx, err := r.db.Begin()
	if err != nil {
		return tracerr.Errorf("failed to begin transaction: %w", err)
	}

	defer tx.Rollback()

	// Check if subscription exists
	var existingSubscription model.Subscription
	checkQuery := `SELECT * FROM subscriptions WHERE user_id = $1`

	switch tx.QueryRow(checkQuery, sub.UserID).Scan(&existingSubscription) {
	case nil:
		if sub.SubscriptionLevel.Priority() < existingSubscription.SubscriptionLevel.Priority() {
			// If the new subscription level is lower than the existing one, do not update
			return appErr.Errorf(appErr.LowPrioritySubscriptionLevelError, "cannot update subscription for user ID %s: new level is lower than existing", sub.UserID)
		}

		// Subscription exists, perform an update
		updateQuery := `
			UPDATE subscriptions
			SET subscription_level = $1, till = $2, since = $3
			WHERE user_id = $4`

		_, err = tx.Exec(updateQuery, sub.SubscriptionLevel, sub.Till, sub.Since, sub.UserID)
		if err != nil {
			return tracerr.Errorf("failed to update subscription for user ID %s: %w", sub.UserID, err)
		}
	case sql.ErrNoRows:
		// Subscription does not exist, perform an insert
		insertQuery := `
			INSERT INTO subscriptions (user_id, subscription_level, till, since)
			VALUES ($1, $2, $3, $4)`

		_, err = tx.Exec(insertQuery, sub.UserID, sub.SubscriptionLevel, sub.Till, sub.Since)
		if err != nil {
			return tracerr.Errorf("failed to insert subscription: %w", err)
		}
	default:
		// Other database error during check
		return tracerr.Errorf("failed to check for existing subscription: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return tracerr.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *pgSubscriptionRepository) GetSubscriptionByUserID(userID uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT user_id, subscription_id, subscription_level, till, since FROM subscriptions
		WHERE user_id = $1`

	sub := &model.Subscription{}
	err := r.db.QueryRow(query, userID).Scan(
		&sub.ID,
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

func (r *pgSubscriptionRepository) RemoveSubscription(id string) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil {
		return tracerr.Errorf("failed to delete subscription%s: %w", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after delete: %w", err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no subscription found to delete %s", id)
	}

	return nil
}
