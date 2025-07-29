package repository

import (
	"context"
	"database/sql"

	appErr "github.com/compendium-tech/compendium/application-assistant-service/internal/error"
	"github.com/compendium-tech/compendium/application-assistant-service/internal/model"
	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgconn"
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

func (r *pgApplicationRepository) CreateApplication(ctx context.Context, app model.Application) error {
	query := `
		INSERT INTO applications (id, user_id, name)
		VALUES ($1, $2, $3)
	`
	_, err := r.db.ExecContext(ctx, query, app.ID, app.UserID, app.Name)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *pgApplicationRepository) GetApplication(ctx context.Context, id uuid.UUID) (*model.Application, error) {
	app := &model.Application{}
	query := `
		SELECT id, user_id, name
		FROM applications
		WHERE id = $1
	`
	row := r.db.QueryRowContext(ctx, query, id)
	err := row.Scan(&app.ID, &app.UserID, &app.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Wrap(err)
	}

	return app, nil
}

func (r *pgApplicationRepository) FindApplicationsByUserID(ctx context.Context, userID uuid.UUID) ([]model.Application, error) {
	query := `
		SELECT id, user_id, name
		FROM applications
		WHERE user_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer rows.Close()

	var applications []model.Application

	for rows.Next() {
		var app model.Application
		if err := rows.Scan(&app.ID, &app.UserID, &app.Name); err != nil {
			return nil, tracerr.Wrap(err)
		}
		applications = append(applications, app)
	}
	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return applications, nil
}

func (r *pgApplicationRepository) RemoveApplication(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM applications
		WHERE id = $1
	`
	res, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return appErr.New(appErr.ApplicationNotFoundError, "application not found")
	}

	return nil
}

func (r *pgApplicationRepository) CreateActivity(ctx context.Context, activity model.Activity) error {
	query := `
		INSERT INTO activities (application_id, name, role, description, hours_per_week, weeks_per_year, category, grades)
		VALUES ($1, $2, $3, $4, $5, $6, $7, $8)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		activity.ApplicationID,
		activity.Name,
		activity.Role,
		activity.Description,
		activity.HoursPerWeek,
		activity.WeeksPerYear,
		activity.Category,
		activity.Grades,
	)
	if err != nil {
		return tracerr.Wrap(err)
	}

	return nil
}

func (r *pgApplicationRepository) GetActivity(ctx context.Context, applicationID uuid.UUID) (*model.Activity, error) {
	activity := &model.Activity{}
	query := `
		SELECT application_id, name, role, description, hours_per_week, weeks_per_year, category, grades
		FROM activities
		WHERE application_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, applicationID)
	err := row.Scan(
		&activity.ApplicationID,
		&activity.Name,
		&activity.Role,
		&activity.Description,
		&activity.HoursPerWeek,
		&activity.WeeksPerYear,
		&activity.Category,
		&activity.Grades,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Wrap(err)
	}

	return activity, nil
}

func (r *pgApplicationRepository) GetActivities(ctx context.Context, applicationID uuid.UUID) ([]model.Activity, error) {
	query := `
		SELECT application_id, name, role, description, hours_per_week, weeks_per_year, category, grades
		FROM activities
		WHERE application_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, applicationID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer rows.Close()

	var activities []model.Activity
	for rows.Next() {
		var activity model.Activity
		if err := rows.Scan(
			&activity.ApplicationID,
			&activity.Name,
			&activity.Role,
			&activity.Description,
			&activity.HoursPerWeek,
			&activity.WeeksPerYear,
			&activity.Category,
			&activity.Grades,
		); err != nil {
			return nil, tracerr.Wrap(err)
		}
		activities = append(activities, activity)
	}

	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return activities, nil
}

func (r *pgApplicationRepository) UpdateActivity(ctx context.Context, activity model.Activity) (*model.Activity, error) {
	query := `
		UPDATE activities
		SET name = $2, role = $3, description = $4, hours_per_week = $5, weeks_per_year = $6, category = $7, grades = $8
		WHERE application_id = $1
		RETURNING application_id, name, role, description, hours_per_week, weeks_per_year, category, grades
	`
	row := r.db.QueryRowContext(
		ctx,
		query,
		activity.ApplicationID,
		activity.Name,
		activity.Role,
		activity.Description,
		activity.HoursPerWeek,
		activity.WeeksPerYear,
		activity.Category,
		activity.Grades,
	)
	updatedActivity := &model.Activity{}
	err := row.Scan(
		&updatedActivity.ApplicationID,
		&updatedActivity.Name,
		&updatedActivity.Role,
		&updatedActivity.Description,
		&updatedActivity.HoursPerWeek,
		&updatedActivity.WeeksPerYear,
		&updatedActivity.Category,
		&updatedActivity.Grades,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.New(appErr.ActivityNotFoundError, "activity not found")
		}

		return nil, tracerr.Wrap(err)
	}
	return updatedActivity, nil
}

func (r *pgApplicationRepository) DeleteActivity(ctx context.Context, applicationID uuid.UUID) error {
	query := `
		DELETE FROM activities
		WHERE application_id = $1
	`
	res, err := r.db.ExecContext(ctx, query, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return appErr.New(appErr.ActivityNotFoundError, "activity not found")
	}

	return nil
}

func (r *pgApplicationRepository) CreateHonor(ctx context.Context, honor model.Honor) error {
	query := `
		INSERT INTO honors (application_id, title, description, level, grade)
		VALUES ($1, $2, $3, $4, $5)
	`
	_, err := r.db.ExecContext(
		ctx,
		query,
		honor.ApplicationID,
		honor.Title,
		honor.Description,
		honor.Level,
		honor.Grade,
	)
	if err != nil {
		if pgErr, ok := err.(*pgconn.PgError); ok && pgErr.Code == "23503" {
			return appErr.New(appErr.ApplicationNotFoundError, "referenced application ID does not exist")
		}

		return tracerr.Wrap(err)
	}
	return nil
}

func (r *pgApplicationRepository) GetHonor(ctx context.Context, applicationID uuid.UUID) (*model.Honor, error) {
	honor := &model.Honor{}
	query := `
		SELECT application_id, title, description, level, grade
		FROM honors
		WHERE application_id = $1
	`
	row := r.db.QueryRowContext(ctx, query, applicationID)
	err := row.Scan(
		&honor.ApplicationID,
		&honor.Title,
		&honor.Description,
		&honor.Level,
		&honor.Grade,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, tracerr.Wrap(err)
	}
	return honor, nil
}

func (r *pgApplicationRepository) GetHonors(ctx context.Context, applicationID uuid.UUID) ([]model.Honor, error) {
	query := `
		SELECT application_id, title, description, level, grade
		FROM honors
		WHERE application_id = $1
	`
	rows, err := r.db.QueryContext(ctx, query, applicationID)
	if err != nil {
		return nil, tracerr.Wrap(err)
	}
	defer rows.Close()

	var honors []model.Honor
	for rows.Next() {
		var honor model.Honor
		if err := rows.Scan(
			&honor.ApplicationID,
			&honor.Title,
			&honor.Description,
			&honor.Level,
			&honor.Grade,
		); err != nil {
			return nil, tracerr.Wrap(err)
		}
		honors = append(honors, honor)
	}
	if err := rows.Err(); err != nil {
		return nil, tracerr.Wrap(err)
	}

	return honors, nil
}

func (r *pgApplicationRepository) UpdateHonor(ctx context.Context, honor model.Honor) (*model.Honor, error) {
	query := `
		UPDATE honors
		SET title = $2, description = $3, level = $4, grade = $5
		WHERE application_id = $1
		RETURNING application_id, title, description, level, grade
	`
	row := r.db.QueryRowContext(
		ctx,
		query,
		honor.ApplicationID,
		honor.Title,
		honor.Description,
		honor.Level,
		honor.Grade,
	)
	updatedHonor := &model.Honor{}
	err := row.Scan(
		&updatedHonor.ApplicationID,
		&updatedHonor.Title,
		&updatedHonor.Description,
		&updatedHonor.Level,
		&updatedHonor.Grade,
	)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, appErr.New(appErr.HonorNotFoundError, "honor not found")
		}

		return nil, tracerr.Wrap(err)
	}
	return updatedHonor, nil
}

func (r *pgApplicationRepository) DeleteHonor(ctx context.Context, applicationID uuid.UUID) error {
	query := `
		DELETE FROM honors
		WHERE application_id = $1
	`
	res, err := r.db.ExecContext(ctx, query, applicationID)
	if err != nil {
		return tracerr.Wrap(err)
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return tracerr.Wrap(err)
	}

	if rowsAffected == 0 {
		return appErr.New(appErr.HonorNotFoundError, "honor not found")
	}

	return nil
}
