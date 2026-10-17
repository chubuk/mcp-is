package storage

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5/pgxpool"
)

type EmailTemplate struct {
	ID           string
	Key          string
	Name         string
	Subject      string
	HTMLBody     string
	TextBody     string
	Description  string
	RequiredVars []string
	IsModified   bool
	UpdatedAt    string
}

type EmailTemplateRepository struct {
	db *pgxpool.Pool
}

func NewEmailTemplateRepository(db *pgxpool.Pool) *EmailTemplateRepository {
	return &EmailTemplateRepository{db: db}
}

const emailTemplateCols = `id, key, name, subject, html_body, COALESCE(text_body,''),
	COALESCE(description,''), COALESCE(required_vars,'{}'), is_modified, updated_at::text`

func (r *EmailTemplateRepository) GetByKey(key string) (*EmailTemplate, error) {
	var t EmailTemplate
	err := r.db.QueryRow(context.Background(),
		`SELECT `+emailTemplateCols+` FROM email_templates WHERE key = $1`, key,
	).Scan(&t.ID, &t.Key, &t.Name, &t.Subject, &t.HTMLBody, &t.TextBody,
		&t.Description, &t.RequiredVars, &t.IsModified, &t.UpdatedAt)
	if err != nil {
		return nil, fmt.Errorf("email template not found: %w", err)
	}
	return &t, nil
}

func (r *EmailTemplateRepository) List() ([]EmailTemplate, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT `+emailTemplateCols+` FROM email_templates ORDER BY key`)
	if err != nil {
		return nil, fmt.Errorf("list email templates: %w", err)
	}
	defer rows.Close()

	var out []EmailTemplate
	for rows.Next() {
		var t EmailTemplate
		if err := rows.Scan(&t.ID, &t.Key, &t.Name, &t.Subject, &t.HTMLBody, &t.TextBody,
			&t.Description, &t.RequiredVars, &t.IsModified, &t.UpdatedAt); err != nil {
			return nil, fmt.Errorf("scan email template: %w", err)
		}
		out = append(out, t)
	}
	return out, nil
}

// Upsert saves a template and marks it modified.
func (r *EmailTemplateRepository) Upsert(key, name, subject, htmlBody, textBody, description string, requiredVars []string) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO email_templates (key, name, subject, html_body, text_body, description, required_vars, is_modified, updated_at)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, TRUE, NOW())
		 ON CONFLICT (key) DO UPDATE SET
		   name = EXCLUDED.name, subject = EXCLUDED.subject,
		   html_body = EXCLUDED.html_body, text_body = EXCLUDED.text_body,
		   description = EXCLUDED.description, required_vars = EXCLUDED.required_vars,
		   is_modified = TRUE, updated_at = NOW()`,
		key, name, subject, htmlBody, textBody, description, requiredVars,
	)
	if err != nil {
		return fmt.Errorf("upsert email template: %w", err)
	}
	return nil
}

// SeedIfMissing inserts defaults without overwriting user edits.
func (r *EmailTemplateRepository) SeedIfMissing(key, name, subject, htmlBody, textBody, description string, requiredVars []string) error {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO email_templates (key, name, subject, html_body, text_body, description, required_vars, is_modified)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, FALSE)
		 ON CONFLICT (key) DO NOTHING`,
		key, name, subject, htmlBody, textBody, description, requiredVars,
	)
	if err != nil {
		return fmt.Errorf("seed email template: %w", err)
	}
	return nil
}

func (r *EmailTemplateRepository) Delete(key string) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM email_templates WHERE key = $1`, key)
	if err != nil {
		return fmt.Errorf("delete email template: %w", err)
	}
	return nil
}
