package organization

import (
	"context"

	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/pubsub"
	"github.com/leg100/otf/internal/sql"
	"github.com/leg100/otf/internal/sql/pggen"
)

type (
	// pgdb is a database of organizations on postgres
	pgdb struct {
		internal.DB // provides access to generated SQL queries
	}

	row struct {
		OrganizationID             pgtype.Text        `json:"organization_id"`
		CreatedAt                  pgtype.Timestamptz `json:"created_at"`
		UpdatedAt                  pgtype.Timestamptz `json:"updated_at"`
		Name                       pgtype.Text        `json:"name"`
		SessionRemember            pgtype.Int4        `json:"session_remember"`
		SessionTimeout             pgtype.Int4        `json:"session_timeout"`
		Email                      pgtype.Text        `json:"email"`
		CollaboratorAuthPolicy     pgtype.Text        `json:"collaborator_auth_policy"`
		AllowForceDeleteWorkspaces bool               `json:"allow_force_delete_workspaces"`
	}

	// dbListOptions represents the options for listing organizations via the
	// database.
	dbListOptions struct {
		names []string // filter organizations by name if non-nil
		internal.ListOptions
	}
)

// GetByID implements pubsub.Getter
func (db *pgdb) GetByID(ctx context.Context, id string, action pubsub.DBAction) (any, error) {
	if action == pubsub.DeleteDBAction {
		return &Organization{ID: id}, nil
	}
	r, err := db.FindOrganizationByID(ctx, sql.String(id))
	if err != nil {
		return nil, sql.Error(err)
	}
	return row(r).toOrganization(), nil
}

func (db *pgdb) update(ctx context.Context, name string, fn func(*Organization) error) (*Organization, error) {
	var org *Organization
	err := db.Tx(ctx, func(tx internal.DB) error {
		result, err := tx.FindOrganizationByNameForUpdate(ctx, sql.String(name))
		if err != nil {
			return err
		}
		org = row(result).toOrganization()

		if err := fn(org); err != nil {
			return err
		}
		_, err = tx.UpdateOrganizationByName(ctx, pggen.UpdateOrganizationByNameParams{
			Name:                       sql.String(name),
			NewName:                    sql.String(org.Name),
			Email:                      sql.StringPtr(org.Email),
			CollaboratorAuthPolicy:     sql.StringPtr(org.CollaboratorAuthPolicy),
			SessionRemember:            sql.Int4Ptr(org.SessionRemember),
			SessionTimeout:             sql.Int4Ptr(org.SessionTimeout),
			UpdatedAt:                  sql.Timestamptz(org.UpdatedAt),
			AllowForceDeleteWorkspaces: org.AllowForceDeleteWorkspaces,
		})
		if err != nil {
			return err
		}
		return nil
	})
	return org, err
}

func (db *pgdb) list(ctx context.Context, opts dbListOptions) (*OrganizationList, error) {
	if opts.names == nil {
		opts.names = []string{"%"} // return all organizations
	}

	batch := &pgx.Batch{}

	db.FindOrganizationsBatch(batch, pggen.FindOrganizationsParams{
		Names:  opts.names,
		Limit:  sql.Int8(opts.GetLimit()),
		Offset: sql.Int8(opts.GetOffset()),
	})
	db.CountOrganizationsBatch(batch, opts.names)
	results := db.SendBatch(ctx, batch)
	defer results.Close()

	rows, err := db.FindOrganizationsScan(results)
	if err != nil {
		return nil, err
	}
	count, err := db.CountOrganizationsScan(results)
	if err != nil {
		return nil, err
	}

	var items []*Organization
	for _, r := range rows {
		items = append(items, row(r).toOrganization())
	}

	return &OrganizationList{
		Items:      items,
		Pagination: internal.NewPagination(opts.ListOptions, int(count.Int)),
	}, nil
}

func (db *pgdb) get(ctx context.Context, name string) (*Organization, error) {
	r, err := db.FindOrganizationByName(ctx, sql.String(name))
	if err != nil {
		return nil, sql.Error(err)
	}
	return row(r).toOrganization(), nil
}

func (db *pgdb) delete(ctx context.Context, name string) error {
	_, err := db.DeleteOrganizationByName(ctx, sql.String(name))
	if err != nil {
		return sql.Error(err)
	}
	return nil
}

// row converts an organization database row into an
// organization.
func (r row) toOrganization() *Organization {
	org := &Organization{
		ID:                         r.OrganizationID.String,
		CreatedAt:                  r.CreatedAt.Time.UTC(),
		UpdatedAt:                  r.UpdatedAt.Time.UTC(),
		Name:                       r.Name.String,
		AllowForceDeleteWorkspaces: r.AllowForceDeleteWorkspaces,
	}
	if r.SessionRemember.Status == pgtype.Present {
		sessionRememberInt := int(r.SessionRemember.Int)
		org.SessionRemember = &sessionRememberInt
	}
	if r.SessionTimeout.Status == pgtype.Present {
		sessionTimeoutInt := int(r.SessionTimeout.Int)
		org.SessionTimeout = &sessionTimeoutInt
	}
	if r.Email.Status == pgtype.Present {
		org.Email = &r.Email.String
	}
	if r.CollaboratorAuthPolicy.Status == pgtype.Present {
		org.CollaboratorAuthPolicy = &r.CollaboratorAuthPolicy.String
	}
	return org
}
