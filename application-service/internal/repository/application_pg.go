package repository

import (
	"context"
	"database/sql"
	"errors"
	"github.com/compendium-tech/compendium/common/pkg/error"
	"github.com/compendium-tech/compendium/common/pkg/pg"

	"github.com/compendium-tech/compendium/application-service/internal/model"
	"github.com/google/uuid"
	"github.com/lib/pq"
	"github.com/ztrue/tracerr"
)

type pgApplicationRepository struct {
	db *sql.DB
}

func NewPgApplicationRepository(db *sql.DB) ApplicationRepository {
	return &pgApplicationRepository{
		db: db,
	}
}

func (r *pgApplicationRepository) GetApplication(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	application := &model.Application{}
	query := `SELECT id, user_id, name, created_at FROM applications WHERE id = $1`
	row := r.db.QueryRowContext(ctx, query, id)

	err := row.Scan(
		&application.ID,
		&application.UserID,
		&application.Name,
		&application.CreatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, nil
		}

		return nil, tracerr.Wrap(err)
	}

	return application, nil
}

func (r *pgApplicationRepository) FindApplicationsByUserID(ctx context.Context, userID uuid.UUID) (_ []model.Application, finalErr error) {
	var applications []model.Application

	query := `SELECT id, user_id, name, created_at FROM applications WHERE user_id = $1`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	defer errorutils.DeferTry(&finalErr, rows.Close)

	for rows.Next() {
		application := model.Application{}
		err := rows.Scan(
			&application.ID,
			&application.UserID,
			&application.Name,
			&application.CreatedAt,
		)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		applications = append(applications, application)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return applications, nil
}

func (r *pgApplicationRepository) CreateApplication(ctx context.Context, app model.Application) error {
	query := `INSERT INTO applications (id, user_id, name, created_at) VALUES ($1, $2, $3, $4)`
	_, err := r.db.ExecContext(
		ctx,
		query,
		app.ID,
		app.UserID,
		app.Name,
		app.CreatedAt,
	)

	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *pgApplicationRepository) UpdateApplicationName(ctx context.Context, applicationID uuid.UUID, name string) error {
	query := `UPDATE applications SET name = $1 WHERE id = $2`
	res, err := r.db.ExecContext(ctx, query, name, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no application found with ID %s to update name", applicationID)
	}

	return nil
}

func (r *pgApplicationRepository) RemoveApplication(ctx context.Context, id uuid.UUID) error {
	query := `DELETE FROM applications WHERE id = $1`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return tracerr.Errorf("no application found with ID %s to remove", id)
	}

	return nil
}

func (r *pgApplicationRepository) GetActivities(ctx context.Context, applicationID uuid.UUID) (_ []model.Activity, finalErr error) {
	var activities []model.Activity
	query := `
		SELECT index, name, role, description, hours_per_week, weeks_per_year, category, grades
		FROM activities
		WHERE application_id = $1
		ORDER BY index
	`
	rows, err := r.db.QueryContext(ctx, query, applicationID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	defer errorutils.DeferTry(&finalErr, rows.Close)

	for rows.Next() {
		activity := model.Activity{}
		var description sql.NullString
		err := rows.Scan(
			&activity.ID,
			&activity.Name,
			&activity.Role,
			&description,
			&activity.HoursPerWeek,
			&activity.WeeksPerYear,
			&activity.Category,
			pq.Array(&activity.Grades),
		)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		if description.Valid {
			activity.Description = &description.String
		}

		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return activities, nil
}

func (r *pgApplicationRepository) PutActivities(ctx context.Context, applicationID uuid.UUID, activities []model.Activity) (finalErr error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return tracerr.Wrap(err)
	}

	defer pg.DeferRollback(&finalErr, tx)

	deleteQuery := `
		DELETE FROM activities
		WHERE application_id = $1
	`
	_, err = tx.ExecContext(ctx, deleteQuery, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	insertQuery := `
		INSERT INTO activities (index, application_id, name, role, description, hours_per_week, weeks_per_year, category, grades)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
	`
	for i, activity := range activities {
		var description sql.NullString
		if activity.Description != nil {
			description = sql.NullString{String: *activity.Description, Valid: true}
		} else {
			description = sql.NullString{Valid: false}
		}

		_, err = tx.ExecContext(
			ctx,
			insertQuery,
			i,
			applicationID,
			activity.Name,
			activity.Role,
			description,
			activity.HoursPerWeek,
			activity.WeeksPerYear,
			activity.Category,
			pq.Array(activity.Grades),
		)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return tx.Commit()
}

func (r *pgApplicationRepository) GetHonors(ctx context.Context, applicationID uuid.UUID) (_ []model.Honor, finalErr error) {
	var honors []model.Honor
	query := `
		SELECT index, application_id, title, description, level, grade
		FROM honors
		WHERE application_id = $1
		ORDER BY index
	`
	rows, err := r.db.QueryContext(ctx, query, applicationID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	defer errorutils.DeferTry(&finalErr, rows.Close)

	for rows.Next() {
		honor := model.Honor{}
		var description sql.NullString
		err := rows.Scan(
			&honor.ID,
			&applicationID,
			&honor.Title,
			&description,
			&honor.Level,
			&honor.Grade,
		)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		if description.Valid {
			honor.Description = &description.String
		}

		honors = append(honors, honor)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return honors, nil
}

func (r *pgApplicationRepository) PutHonors(ctx context.Context, applicationID uuid.UUID, honors []model.Honor) (finalErr error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return tracerr.Wrap(err)
	}

	defer pg.DeferRollback(&finalErr, tx)

	deleteQuery := `
		DELETE FROM honors
		WHERE application_id = $1
	`
	_, err = tx.ExecContext(ctx, deleteQuery, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	insertQuery := `
		INSERT INTO honors (index, application_id, title, description, level, grade)
		VALUES ($1, $2, $3, $4, $5, $6)
	`
	for i, honor := range honors {
		var description sql.NullString
		if honor.Description != nil {
			description = sql.NullString{String: *honor.Description, Valid: true}
		} else {
			description = sql.NullString{Valid: false}
		}

		_, err = tx.ExecContext(
			ctx,
			insertQuery,
			i,
			applicationID,
			honor.Title,
			description,
			honor.Level,
			honor.Grade,
		)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return tx.Commit()
}

func (r *pgApplicationRepository) GetEssays(ctx context.Context, applicationID uuid.UUID) (_ []model.Essay, finalErr error) {
	var essays []model.Essay
	query := `
		SELECT index, application_id, type, content
		FROM essays
		WHERE application_id = $1
		ORDER BY index
	`
	rows, err := r.db.QueryContext(ctx, query, applicationID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	defer errorutils.DeferTry(&finalErr, rows.Close)

	for rows.Next() {
		essay := model.Essay{}
		err := rows.Scan(
			&essay.ID,
			&applicationID,
			&essay.Type,
			&essay.Content,
		)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		essays = append(essays, essay)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return essays, nil
}

func (r *pgApplicationRepository) PutEssays(ctx context.Context, applicationID uuid.UUID, essays []model.Essay) (finalErr error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return tracerr.Wrap(err)
	}

	defer pg.DeferRollback(&finalErr, tx)

	deleteQuery := `
		DELETE FROM essays
		WHERE application_id = $1
	`
	_, err = tx.ExecContext(ctx, deleteQuery, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	insertQuery := `
		INSERT INTO essays (index, application_id, type, content)
		VALUES ($1, $2, $3, $4)
	`
	for i, essay := range essays {
		_, err = tx.ExecContext(
			ctx,
			insertQuery,
			i,
			applicationID,
			essay.Type,
			essay.Content,
		)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return tx.Commit()
}

func (r *pgApplicationRepository) GetSupplementalEssays(ctx context.Context, applicationID uuid.UUID) (_ []model.SupplementalEssay, finalErr error) {
	var supplementalEssays []model.SupplementalEssay
	query := `
		SELECT index, title, content FROM supplemental_essays
		WHERE application_id = $1 ORDER BY index`
	rows, err := r.db.QueryContext(ctx, query, applicationID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}

	defer errorutils.DeferTry(&finalErr, rows.Close)

	for rows.Next() {
		supplementalEssay := model.SupplementalEssay{}
		err := rows.Scan(
			&supplementalEssay.ID,
			&supplementalEssay.Title,
			&supplementalEssay.Content,
		)
		if err != nil {
			return nil, tracerr.Wrap(err)
		}

		supplementalEssays = append(supplementalEssays, supplementalEssay)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return supplementalEssays, nil
}

func (r *pgApplicationRepository) PutSupplementalEssays(ctx context.Context, applicationID uuid.UUID, supplementalEssays []model.SupplementalEssay) (finalErr error) {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return tracerr.Wrap(err)
	}

	defer pg.DeferRollback(&finalErr, tx)

	deleteQuery := `
		DELETE FROM supplemental_essays
		WHERE application_id = $1
	`
	_, err = tx.ExecContext(ctx, deleteQuery, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	insertQuery := `
		INSERT INTO supplemental_essays (index, application_id, title, content)
		VALUES ($1, $2, $3, $4)
	`
	for i, supplementalEssay := range supplementalEssays {
		_, err = tx.ExecContext(
			ctx,
			insertQuery,
			i,
			applicationID,
			supplementalEssay.Title,
			supplementalEssay.Content,
		)
		if err != nil {
			return tracerr.Wrap(err)
		}
	}

	return tx.Commit()
}
