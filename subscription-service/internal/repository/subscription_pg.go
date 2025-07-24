package repository

import (
	"context"
	"database/sql"
	"time"

	appErr "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/ztrue/tracerr"
)

type pgSubscriptionRepository struct {
	db *sql.DB
}

func NewPgSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &pgSubscriptionRepository{db: db}
}

func (r *pgSubscriptionRepository) PutSubscription(ctx context.Context, sub model.Subscription) error {
	tx, err := r.db.Begin()
	if err != nil {
		return tracerr.Errorf("failed to begin transaction: %w", err)
	}

	defer func() {
		if r := recover(); r != nil {
			tx.Rollback()
			panic(r) // Re-throw panic
		} else if err != nil {
			tx.Rollback()
		}
	}()

	// Check if subscription exists
	var existingSubscription model.Subscription
	// Note: We need to scan the invitation_code as well now
	checkQuery := `SELECT id, backed_by, subscription_level, till, since, invitation_code FROM subscriptions WHERE backed_by = $1`

	var inviteCodeSQL sql.NullString // Use sql.NullString to handle nullable invitation_code
	err = tx.QueryRowContext(ctx, checkQuery, sub.BackedBy).Scan(&existingSubscription.ID, &existingSubscription.BackedBy, &existingSubscription.Level, &existingSubscription.Till, &existingSubscription.Since, &inviteCodeSQL)

	switch err {
	case nil:
		if inviteCodeSQL.Valid {
			existingSubscription.InvitationCode = &inviteCodeSQL.String
		} else {
			existingSubscription.InvitationCode = nil
		}

		if sub.Level.Priority() < existingSubscription.Level.Priority() {
			// If the new subscription level is lower than the existing one, do not update
			return appErr.Errorf(appErr.LowPrioritySubscriptionLevelError, "cannot update subscription for user ID %s: new level is lower than existing", sub.BackedBy)
		}

		// Subscription exists, perform an update
		// Include invitation_code in update if it's explicitly provided
		updateQuery := `
			UPDATE subscriptions
			SET subscription_level = $1, till = $2, since = $3, invitation_code = $4
			WHERE backed_by = $5`

		_, err = tx.ExecContext(ctx, updateQuery, sub.Level, sub.Till, sub.Since, sub.InvitationCode, sub.BackedBy)
		if err != nil {
			return tracerr.Errorf("failed to update subscription for user ID %s: %w", sub.BackedBy, err)
		}
	case sql.ErrNoRows:
		// Subscription does not exist, perform an insert
		insertQuery := `
			INSERT INTO subscriptions (backed_by, subscription_level, till, since, invitation_code)
			VALUES ($1, $2, $3, $4, $5)`

		_, err = tx.ExecContext(ctx, insertQuery, sub.BackedBy, sub.Level, sub.Till, sub.Since, sub.InvitationCode)
		if err != nil {
			return tracerr.Errorf("failed to insert subscription: %w", err)
		}

		// Insert into subscription_members table
		insertMemberQuery := `
			INSERT INTO subscription_members (subscription_id, user_id, since)
			VALUES ($1, $2, $3)`

		_, err = tx.ExecContext(ctx, insertMemberQuery, sub.ID, sub.BackedBy, time.Now().UTC())
		if err != nil {
			return tracerr.Errorf("failed to insert subscription member: %w", err)
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

func (r *pgSubscriptionRepository) GetSubscriptionByInvitationCode(ctx context.Context, invitationCode string) (*model.Subscription, error) {
	query := `
		SELECT id, backed_by, subscription_level, till, since
		FROM subscriptions
		WHERE invitation_code = $1`

	var sub model.Subscription
	err := r.db.QueryRowContext(ctx, query, invitationCode).Scan(&sub.ID, &sub.BackedBy, &sub.Level, &sub.Till, &sub.Since)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Errorf("failed to get subscription by invitation code %s: %w", invitationCode, err)
	}

	sub.InvitationCode = &invitationCode

	return &sub, nil
}

func (r *pgSubscriptionRepository) GetSubscriptionByMemberUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT s.id, s.backed_by, s.subscription_level, s.till, s.since, s.invitation_code
		FROM subscriptions s
		INNER JOIN subscription_members sm ON s.id = sm.subscription_id
		WHERE sm.user_id = $1`
	var sub model.Subscription
	var inviteCodeSQL sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sub.ID, &sub.BackedBy, &sub.Level, &sub.Till, &sub.Since, &inviteCodeSQL)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Errorf("failed to get subscription for user ID %s: %w", userID, err)
	}

	if inviteCodeSQL.Valid {
		sub.InvitationCode = &inviteCodeSQL.String
	} else {
		sub.InvitationCode = nil
	}

	return &sub, nil
}

func (r *pgSubscriptionRepository) GetSubscriptionByPayerUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT id, backed_by, subscription_level, till, since, invitation_code
		FROM subscriptions
		WHERE backed_by = $1`

	var sub model.Subscription
	var inviteCodeSQL sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sub.ID, &sub.BackedBy, &sub.Level, &sub.Till, &sub.Since, &inviteCodeSQL)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Errorf("failed to get subscription for user ID %s: %w", userID, err)
	}

	if inviteCodeSQL.Valid {
		sub.InvitationCode = &inviteCodeSQL.String
	} else {
		sub.InvitationCode = nil
	}

	return &sub, nil
}

func (r *pgSubscriptionRepository) GetSubscriptionMembers(ctx context.Context, subscriptionID string) ([]model.SubscriptionMember, error) {
	query := `
		SELECT user_id, since
		FROM subscription_members
		WHERE subscription_id = $1`

	rows, err := r.db.QueryContext(ctx, query, subscriptionID)
	if err != nil {
		return nil, tracerr.Errorf("failed to get subscription members for ID %s: %w", subscriptionID, err)
	}
	defer rows.Close()

	var members []model.SubscriptionMember
	for rows.Next() {
		var member model.SubscriptionMember
		if err := rows.Scan(&member.UserID, &member.Since); err != nil {
			return nil, tracerr.Errorf("failed to scan subscription member: %w", err)
		}

		member.SubscriptionID = subscriptionID
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Errorf("error occurred while iterating over subscription members: %w", err)
	}

	return members, nil
}

func (r *pgSubscriptionRepository) AddSubscriptionMember(ctx context.Context, member model.SubscriptionMember) error {
	query := `
		INSERT INTO subscription_members (subscription_id, user_id, since)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, member.SubscriptionID, member.UserID, member.Since)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23505" {
			return appErr.Errorf(appErr.AlreadySubscribedError, "you are already member of subscription: %s", member.SubscriptionID)
		}

		return tracerr.Errorf("failed to add subscription member: %w", err)
	}

	return nil
}

// UpdateSubscriptionInvitationCode updates the invitation code for a given subscription.
// If code is nil, it sets the invitation_code to NULL in the database.
func (r *pgSubscriptionRepository) UpdateSubscriptionInvitationCode(ctx context.Context, subscriptionID string, code *string) error {
	query := `
		UPDATE subscriptions
		SET invitation_code = $1
		WHERE id = $2`

	var sqlCode sql.NullString
	if code != nil {
		sqlCode.String = *code
		sqlCode.Valid = true
	} else {
		sqlCode.Valid = false // This will set the column to NULL
	}

	result, err := r.db.ExecContext(ctx, query, sqlCode, subscriptionID)
	if err != nil {
		return tracerr.Errorf("failed to update invitation code for subscription ID %s: %w", subscriptionID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after updating invitation code: %w", err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no subscription found with ID %s to update invitation code", subscriptionID)
	}

	return nil
}

func (r *pgSubscriptionRepository) RemoveSubscriptionMember(ctx context.Context, member model.SubscriptionMember) error {
	query := `
		DELETE FROM subscription_members
		WHERE subscription_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, member.SubscriptionID, member.UserID)
	if err != nil {
		return tracerr.Errorf("failed to delete subscription member: %w", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after delete: %w", err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no subscription member found to delete for user ID %s in subscription %s", member.UserID, member.SubscriptionID)
	}

	return nil
}

func (r *pgSubscriptionRepository) RemoveSubscriptionMemberBySubscriptionAndUserID(ctx context.Context, subscriptionID string, userID uuid.UUID) error {
	query := `
		DELETE FROM subscription_members
		WHERE subscription_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, subscriptionID, userID)
	if err != nil {
		return tracerr.Errorf("failed to delete subscription member for user ID %s in subscription %s: %w", userID, subscriptionID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after delete: %w", err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no subscription member found to delete for user ID %s in subscription %s", userID, subscriptionID)
	}

	return nil
}

func (r *pgSubscriptionRepository) RemoveSubscription(ctx context.Context, id string) error {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
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
