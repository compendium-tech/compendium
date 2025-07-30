package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/compendium-tech/compendium/common/pkg/error"
	"time"

	"github.com/compendium-tech/compendium/common/pkg/log"
	"github.com/compendium-tech/compendium/common/pkg/pg"
	"github.com/compendium-tech/compendium/subscription-service/internal/error"
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

func (r *pgSubscriptionRepository) UpsertSubscription(ctx context.Context, sub model.Subscription) (finalErr error) {
	tx, err := r.db.Begin()
	if err != nil {
		return tracerr.Errorf("failed to begin transaction: %v", err)
	}

	defer pg.DeferRollback(&finalErr, tx)

	// Check if the user is a subscription payer or member
	var existingSubscriptionID string
	var existingSubscriptionBackedBy uuid.UUID
	var inviteCodeSQL sql.NullString
	query := `
        SELECT s.id, s.backed_by
        FROM subscriptions s
        JOIN subscription_members sm ON s.id = sm.subscription_id
        WHERE sm.user_id = $1`

	log.L(ctx).Infof("Checking for existing subscription or membership for user ID %s", sub.BackedBy)
	err = tx.QueryRowContext(ctx, query, sub.BackedBy).Scan(
		&existingSubscriptionID,
		&existingSubscriptionBackedBy,
		&inviteCodeSQL,
	)

	if err == nil {
		log.L(ctx).Infof("Found existing subscription or membership for user ID %s", sub.BackedBy)

		// Check if the user is the payer (backed_by)
		if existingSubscriptionBackedBy == sub.BackedBy {
			// User is the payer, delete the subscription (cascade will remove members)
			log.L(ctx).Infof("User is subscription payer, deleting subscription %s", existingSubscriptionID)
			deleteQuery := `DELETE FROM subscriptions WHERE id = $1`
			result, err := tx.ExecContext(ctx, deleteQuery, existingSubscriptionID)
			if err != nil {
				return tracerr.Errorf("failed to delete existing subscription for user ID %s: %v", sub.BackedBy, err)
			}

			// Verify at least one row was deleted
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return tracerr.Errorf("failed to check rows affected: %v", err)
			}
			if rowsAffected == 0 {
				return tracerr.Errorf("no subscription was deleted for user ID %s", sub.BackedBy)
			}

			log.L(ctx).Infof("Successfully deleted subscription for user ID %s", sub.BackedBy)
		} else {
			// User is just a member, delete only their membership

			log.L(ctx).Infof("User is subscription member, deleting membership for subscription %s", existingSubscriptionID)
			deleteMemberQuery := `DELETE FROM subscription_members WHERE subscription_id = $1 AND user_id = $2`
			result, err := tx.ExecContext(ctx, deleteMemberQuery, existingSubscriptionID, sub.BackedBy)
			if err != nil {
				return tracerr.Errorf("failed to delete subscription member for user ID %s: %v", sub.BackedBy, err)
			}

			// Verify at least one row was deleted
			rowsAffected, err := result.RowsAffected()
			if err != nil {
				return tracerr.Errorf("failed to check rows affected: %v", err)
			}
			if rowsAffected == 0 {
				return tracerr.Errorf("no subscription member was deleted for user ID %s", sub.BackedBy)
			}

			log.L(ctx).Infof("Successfully deleted subscription member for user ID %s", sub.BackedBy)
		}
	} else if !errors.Is(err, sql.ErrNoRows) {
		// Other database error during check
		return tracerr.Errorf("failed to check for existing subscription or membership: %v", err)
	} else {
		log.L(ctx).Infof("No existing subscription or membership found for user ID %s", sub.BackedBy)
	}

	// Insert new subscription
	log.L(ctx).Infof("Inserting new subscription for user ID %s", sub.BackedBy)
	insertQuery := `
        INSERT INTO subscriptions (id, backed_by, tier, till, since, invitation_code)
        VALUES ($1, $2, $3, $4, $5, $6)`

	_, err = tx.ExecContext(ctx, insertQuery, sub.ID, sub.BackedBy, sub.Tier, sub.Till, sub.Since, sub.InvitationCode)
	if err != nil {
		return tracerr.Errorf("failed to insert subscription: %w", err)
	}

	log.L(ctx).Infof("Successfully inserted subscription for user ID %s", sub.BackedBy)

	// Insert into subscription_members table
	log.L(ctx).Infof("Inserting subscription member for user ID %s", sub.BackedBy)
	insertMemberQuery := `
        INSERT INTO subscription_members (subscription_id, user_id, since)
        VALUES ($1, $2, $3)`

	_, err = tx.ExecContext(ctx, insertMemberQuery, sub.ID, sub.BackedBy, time.Now().UTC())
	if err != nil {
		return tracerr.Errorf("failed to insert subscription member: %w", err)
	}

	err = tx.Commit()
	if err != nil {
		return tracerr.Errorf("failed to commit transaction: %w", err)
	}

	return nil
}

func (r *pgSubscriptionRepository) FindSubscriptionByInvitationCode(ctx context.Context, invitationCode string) (*model.Subscription, error) {
	query := `
		SELECT id, backed_by, tier, till, since
		FROM subscriptions
		WHERE invitation_code = $1`

	var sub model.Subscription
	err := r.db.QueryRowContext(ctx, query, invitationCode).Scan(&sub.ID, &sub.BackedBy, &sub.Tier, &sub.Till, &sub.Since)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}
		return nil, tracerr.Errorf("failed to get subscription by invitation code %s: %w", invitationCode, err)
	}

	sub.InvitationCode = &invitationCode

	return &sub, nil
}

func (r *pgSubscriptionRepository) FindSubscriptionByMemberUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
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

func (r *pgSubscriptionRepository) FindSubscriptionByPayerUserID(ctx context.Context, userID uuid.UUID) (*model.Subscription, error) {
	query := `
		SELECT id, backed_by, tier, till, since, invitation_code
		FROM subscriptions
		WHERE backed_by = $1`

	var sub model.Subscription
	var inviteCodeSQL sql.NullString
	err := r.db.QueryRowContext(ctx, query, userID).Scan(&sub.ID, &sub.BackedBy, &sub.Tier, &sub.Till, &sub.Since, &inviteCodeSQL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
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

func (r *pgSubscriptionRepository) GetSubscriptionMembers(ctx context.Context, subscriptionID string) (_ []model.SubscriptionMember, finalErr error) {
	query := `
		SELECT user_id, since
		FROM subscription_members
		WHERE subscription_id = $1`

	rows, err := r.db.QueryContext(ctx, query, subscriptionID)
	if err != nil {
		return nil, tracerr.Errorf("failed to get subscription members for ID %s: %w", subscriptionID, err)
	}

	defer errorutils.DeferTry(&finalErr, rows.Close)

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

func (r *pgSubscriptionRepository) CreateSubscriptionMember(ctx context.Context, member model.SubscriptionMember) error {
	query := `
		INSERT INTO subscription_members (subscription_id, user_id, since)
		VALUES ($1, $2, $3)`

	_, err := r.db.ExecContext(ctx, query, member.SubscriptionID, member.UserID, member.Since)
	if err != nil {
		var pgErr *pgconn.PgError
		if errors.As(err, &pgErr) && pgErr.Code == "23505" {
			return myerror.New(myerror.AlreadySubscribedError)
		}

		return tracerr.Errorf("failed to add subscription member: %v", err)
	}

	return nil
}

func (r *pgSubscriptionRepository) RemoveSubscriptionMember(ctx context.Context, member model.SubscriptionMember) error {
	query := `
		DELETE FROM subscription_members
		WHERE subscription_id = $1 AND user_id = $2`
	result, err := r.db.ExecContext(ctx, query, member.SubscriptionID, member.UserID)
	if err != nil {
		return tracerr.Errorf("failed to delete subscription member: %v", err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after delete: %v", err)
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
		return tracerr.Errorf("failed to delete subscription member for user ID %s in subscription %s: %v", userID, subscriptionID, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after delete: %v", err)
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
		return tracerr.Errorf("failed to delete subscription %s: %v", id, err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return tracerr.Errorf("failed to get rows affected after delete: %v", err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no subscription found to delete %s", id)
	}

	return nil
}
