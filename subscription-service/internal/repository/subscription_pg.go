package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"

	"github.com/compendium-tech/compendium/common/pkg/log"
	myerror "github.com/compendium-tech/compendium/subscription-service/internal/error"
	"github.com/compendium-tech/compendium/subscription-service/internal/model"
)

type pgSubscriptionRepository struct {
	db *sql.DB
}

func NewPgSubscriptionRepository(db *sql.DB) SubscriptionRepository {
	return &pgSubscriptionRepository{db: db}
}

func (r *pgSubscriptionRepository) UpsertSubscription(ctx context.Context, sub model.Subscription) {
	tx, err := r.db.Begin()
	if err != nil {
		panic(fmt.Errorf("failed to begin transaction: %v", err))
	}

	defer tx.Rollback()

	// Check for an existing subscription with the same ID.
	var existingID string
	query := `SELECT id FROM subscriptions WHERE id = $1`

	err = tx.QueryRowContext(ctx, query, sub.ID).Scan(&existingID)

	if err == nil {
		// A subscription with this ID already exists. Update it and its membership.
		log.L(ctx).Infof("Found existing subscription with ID %s. Updating it.", sub.ID)

		// Update the existing subscription's details.
		updateQuery := `
            UPDATE subscriptions
            SET backed_by = $2, tier = $3, till = $4, since = $5, invitation_code = $6
            WHERE id = $1`

		_, err = tx.ExecContext(ctx, updateQuery, sub.ID, sub.BackedBy, sub.Tier, sub.Till, sub.Since, sub.InvitationCode)
		if err != nil {
			panic(fmt.Errorf("failed to update subscription with ID %s: %w", sub.ID, err))
		}

		// Now, update the membership record for the payer.
		// We use UPSERT logic here to handle cases where the payer might have been removed.
		upsertMemberQuery := `
            INSERT INTO subscription_members (subscription_id, user_id, since)
            VALUES ($1, $2, $3)
            ON CONFLICT (subscription_id, user_id) DO UPDATE SET since = EXCLUDED.since`

		_, err = tx.ExecContext(ctx, upsertMemberQuery, sub.ID, sub.BackedBy, time.Now().UTC())
		if err != nil {
			panic(fmt.Errorf("failed to upsert subscription member: %w", err))
		}

	} else if errors.Is(err, sql.ErrNoRows) {
		// No subscription with this ID exists. Insert a new one.
		log.L(ctx).Infof("No existing subscription found with ID %s. Inserting a new one.", sub.ID)

		// Insert new subscription
		insertSubscriptionQuery := `
            INSERT INTO subscriptions (id, backed_by, tier, till, since, invitation_code)
            VALUES ($1, $2, $3, $4, $5, $6)`

		_, err = tx.ExecContext(ctx, insertSubscriptionQuery, sub.ID, sub.BackedBy, sub.Tier, sub.Till, sub.Since, sub.InvitationCode)
		if err != nil {
			panic(fmt.Errorf("failed to insert subscription: %w", err))
		}

		// Insert the payer into the subscription_members table.
		insertMemberQuery := `
            INSERT INTO subscription_members (subscription_id, user_id, since)
            VALUES ($1, $2, $3)`

		_, err = tx.ExecContext(ctx, insertMemberQuery, sub.ID, sub.BackedBy, time.Now().UTC())
		if err != nil {
			panic(fmt.Errorf("failed to insert subscription member: %w", err))
		}
	} else {
		panic(fmt.Errorf("failed to check for existing subscription: %v", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("failed to commit transaction: %w", err))
	}

	log.L(ctx).Infof("Successfully completed UpsertSubscription for ID %s", sub.ID)
}

func (r *pgSubscriptionRepository) FindSubscriptionByInvitationCode(ctx context.Context, invitationCode string) *model.Subscription {
	query := `
		SELECT id, backed_by, tier, till, since
		FROM subscriptions
		WHERE invitation_code = $1`

	var sub model.Subscription
	err := r.db.QueryRowContext(ctx, query, invitationCode).Scan(&sub.ID, &sub.BackedBy, &sub.Tier, &sub.Till, &sub.Since)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		panic(fmt.Errorf("failed to get subscription by invitation code %s: %w", invitationCode, err))
	}

	sub.InvitationCode = &invitationCode

	return &sub
}

func (r *pgSubscriptionRepository) FindSubscriptionByMemberUserID(ctx context.Context, userID uuid.UUID) *model.Subscription {
	query := `
		SELECT s.id, s.backed_by, s.tier, s.till, s.since, s.invitation_code
		FROM subscriptions s
		JOIN subscription_members sm ON s.id = sm.subscription_id
		WHERE sm.user_id = $1`
	var sub model.Subscription
	var inviteCodeSQL sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sub.ID, &sub.BackedBy, &sub.Tier, &sub.Till, &sub.Since, &inviteCodeSQL)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		panic(fmt.Errorf("failed to get subscription for user ID %s: %w", userID, err))
	}

	if inviteCodeSQL.Valid {
		sub.InvitationCode = &inviteCodeSQL.String
	} else {
		sub.InvitationCode = nil
	}

	return &sub
}

func (r *pgSubscriptionRepository) FindSubscriptionByPayerUserID(ctx context.Context, userID uuid.UUID) *model.Subscription {
	query := `
		SELECT id, backed_by, tier, till, since, invitation_code
		FROM subscriptions
		WHERE backed_by = $1`

	var sub model.Subscription
	var inviteCodeSQL sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sub.ID, &sub.BackedBy, &sub.Tier, &sub.Till, &sub.Since, &inviteCodeSQL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil
		}

		panic(fmt.Errorf("failed to get subscription for user ID %s: %w", userID, err))
	}

	if inviteCodeSQL.Valid {
		sub.InvitationCode = &inviteCodeSQL.String
	} else {
		sub.InvitationCode = nil
	}

	return &sub
}

func (r *pgSubscriptionRepository) GetSubscriptionMembers(ctx context.Context, subscriptionID string) []model.SubscriptionMember {
	query := `
		SELECT user_id, since
		FROM subscription_members
		WHERE subscription_id = $1`

	rows, err := r.db.QueryContext(ctx, query, subscriptionID)
	if err != nil {
		panic(fmt.Errorf("failed to get subscription members for ID %s: %w", subscriptionID, err))
	}

	defer rows.Close()

	var members []model.SubscriptionMember
	for rows.Next() {
		var member model.SubscriptionMember
		if err := rows.Scan(&member.UserID, &member.Since); err != nil {
			panic(fmt.Errorf("failed to scan subscription member: %w", err))
		}

		member.SubscriptionID = subscriptionID
		members = append(members, member)
	}

	if err := rows.Err(); err != nil {
		panic(fmt.Errorf("error occurred while iterating over subscription members: %w", err))
	}

	return members
}

func (r *pgSubscriptionRepository) CreateSubscriptionMemberAndCheckMemberCount(
	ctx context.Context,
	member model.SubscriptionMember,
	checkCount func(uint) error,
) {
	tx, err := r.db.Begin()
	if err != nil {
		panic(fmt.Errorf("failed to begin transaction: %w", err))
	}

	defer tx.Rollback()

	var currentMembers uint
	queryCountAndLock := `
		SELECT count(sm.user_id)
		FROM subscription_members AS sm
		INNER JOIN subscriptions AS s ON s.id = sm.subscription_id
		WHERE s.id = $1
		FOR UPDATE OF s, sm`

	row := tx.QueryRowContext(ctx, queryCountAndLock, member.SubscriptionID)
	if err := row.Scan(&currentMembers); err != nil {
		panic(fmt.Errorf("failed to get member count and lock subscription: %w", err))
	}

	if err := checkCount(currentMembers); err != nil {
		panic(err)
	}

	queryInsert := `
		INSERT INTO subscription_members (subscription_id, user_id, since)
		VALUES ($1, $2, $3)`

	_, err = tx.ExecContext(ctx, queryInsert, member.SubscriptionID, member.UserID, member.Since)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			myerror.New(myerror.AlreadySubscribedError).Throw()
		}

		panic(fmt.Errorf("failed to add subscription member: %w", err))
	}

	err = tx.Commit()
	if err != nil {
		panic(fmt.Errorf("failed to commit transaction: %w", err))
	}
}

func (r *pgSubscriptionRepository) RemoveSubscriptionMember(ctx context.Context, member model.SubscriptionMember) {
	query := `
		DELETE FROM subscription_members
		WHERE subscription_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, member.SubscriptionID, member.UserID)
	if err != nil {
		panic(fmt.Errorf("failed to delete subscription member: %w", err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Errorf("failed to get rows affected after delete: %w", err))
	}

	if rowsAffected == 0 {
		panic(fmt.Errorf("no subscription member found to delete for user ID %s in subscription %s", member.UserID, member.SubscriptionID))
	}
}

func (r *pgSubscriptionRepository) RemoveSubscriptionMemberBySubscriptionAndUserID(ctx context.Context, subscriptionID string, userID uuid.UUID) {
	query := `
		DELETE FROM subscription_members
		WHERE subscription_id = $1 AND user_id = $2`

	result, err := r.db.ExecContext(ctx, query, subscriptionID, userID)
	if err != nil {
		panic(fmt.Errorf("failed to delete subscription member for user ID %s in subscription %s: %w", userID, subscriptionID, err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Errorf("failed to get rows affected after delete: %w", err))
	}

	if rowsAffected == 0 {
		panic(fmt.Errorf("no subscription member found to delete for user ID %s in subscription %s", userID, subscriptionID))
	}
}

func (r *pgSubscriptionRepository) RemoveSubscription(ctx context.Context, id string) {
	query := `
		DELETE FROM subscriptions
		WHERE id = $1`

	result, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		panic(fmt.Errorf("failed to delete subscription %s: %w", id, err))
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		panic(fmt.Errorf("failed to get rows affected after delete: %w", err))
	}

	if rowsAffected == 0 {
		panic(fmt.Errorf("no subscription found to delete %s", id))
	}
}
