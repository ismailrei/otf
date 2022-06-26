// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertPlanSQL = `INSERT INTO plans (
    run_id,
    status
) VALUES (
    $1,
    $2
);`

// InsertPlan implements Querier.InsertPlan.
func (q *DBQuerier) InsertPlan(ctx context.Context, runID pgtype.Text, status pgtype.Text) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertPlan")
	cmdTag, err := q.conn.Exec(ctx, insertPlanSQL, runID, status)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertPlan: %w", err)
	}
	return cmdTag, err
}

// InsertPlanBatch implements Querier.InsertPlanBatch.
func (q *DBQuerier) InsertPlanBatch(batch genericBatch, runID pgtype.Text, status pgtype.Text) {
	batch.Queue(insertPlanSQL, runID, status)
}

// InsertPlanScan implements Querier.InsertPlanScan.
func (q *DBQuerier) InsertPlanScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertPlanBatch: %w", err)
	}
	return cmdTag, err
}

const updatePlanStatusByIDSQL = `UPDATE plans
SET status = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdatePlanStatusByID implements Querier.UpdatePlanStatusByID.
func (q *DBQuerier) UpdatePlanStatusByID(ctx context.Context, status pgtype.Text, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanStatusByID")
	row := q.conn.QueryRow(ctx, updatePlanStatusByIDSQL, status, runID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdatePlanStatusByID: %w", err)
	}
	return item, nil
}

// UpdatePlanStatusByIDBatch implements Querier.UpdatePlanStatusByIDBatch.
func (q *DBQuerier) UpdatePlanStatusByIDBatch(batch genericBatch, status pgtype.Text, runID pgtype.Text) {
	batch.Queue(updatePlanStatusByIDSQL, status, runID)
}

// UpdatePlanStatusByIDScan implements Querier.UpdatePlanStatusByIDScan.
func (q *DBQuerier) UpdatePlanStatusByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdatePlanStatusByIDBatch row: %w", err)
	}
	return item, nil
}

const updatePlannedChangesByIDSQL = `UPDATE plans
SET report = (
    $1,
    $2,
    $3
)
WHERE run_id = $4
RETURNING run_id
;`

type UpdatePlannedChangesByIDParams struct {
	Additions    int
	Changes      int
	Destructions int
	RunID        pgtype.Text
}

// UpdatePlannedChangesByID implements Querier.UpdatePlannedChangesByID.
func (q *DBQuerier) UpdatePlannedChangesByID(ctx context.Context, params UpdatePlannedChangesByIDParams) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlannedChangesByID")
	row := q.conn.QueryRow(ctx, updatePlannedChangesByIDSQL, params.Additions, params.Changes, params.Destructions, params.RunID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdatePlannedChangesByID: %w", err)
	}
	return item, nil
}

// UpdatePlannedChangesByIDBatch implements Querier.UpdatePlannedChangesByIDBatch.
func (q *DBQuerier) UpdatePlannedChangesByIDBatch(batch genericBatch, params UpdatePlannedChangesByIDParams) {
	batch.Queue(updatePlannedChangesByIDSQL, params.Additions, params.Changes, params.Destructions, params.RunID)
}

// UpdatePlannedChangesByIDScan implements Querier.UpdatePlannedChangesByIDScan.
func (q *DBQuerier) UpdatePlannedChangesByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdatePlannedChangesByIDBatch row: %w", err)
	}
	return item, nil
}

const getPlanBinByIDSQL = `SELECT plan_bin
FROM plans
WHERE run_id = $1
;`

// GetPlanBinByID implements Querier.GetPlanBinByID.
func (q *DBQuerier) GetPlanBinByID(ctx context.Context, runID pgtype.Text) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetPlanBinByID")
	row := q.conn.QueryRow(ctx, getPlanBinByIDSQL, runID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetPlanBinByID: %w", err)
	}
	return item, nil
}

// GetPlanBinByIDBatch implements Querier.GetPlanBinByIDBatch.
func (q *DBQuerier) GetPlanBinByIDBatch(batch genericBatch, runID pgtype.Text) {
	batch.Queue(getPlanBinByIDSQL, runID)
}

// GetPlanBinByIDScan implements Querier.GetPlanBinByIDScan.
func (q *DBQuerier) GetPlanBinByIDScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetPlanBinByIDBatch row: %w", err)
	}
	return item, nil
}

const getPlanJSONByIDSQL = `SELECT plan_json
FROM plans
WHERE run_id = $1
;`

// GetPlanJSONByID implements Querier.GetPlanJSONByID.
func (q *DBQuerier) GetPlanJSONByID(ctx context.Context, runID pgtype.Text) ([]byte, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "GetPlanJSONByID")
	row := q.conn.QueryRow(ctx, getPlanJSONByIDSQL, runID)
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query GetPlanJSONByID: %w", err)
	}
	return item, nil
}

// GetPlanJSONByIDBatch implements Querier.GetPlanJSONByIDBatch.
func (q *DBQuerier) GetPlanJSONByIDBatch(batch genericBatch, runID pgtype.Text) {
	batch.Queue(getPlanJSONByIDSQL, runID)
}

// GetPlanJSONByIDScan implements Querier.GetPlanJSONByIDScan.
func (q *DBQuerier) GetPlanJSONByIDScan(results pgx.BatchResults) ([]byte, error) {
	row := results.QueryRow()
	item := []byte{}
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan GetPlanJSONByIDBatch row: %w", err)
	}
	return item, nil
}

const updatePlanBinByIDSQL = `UPDATE plans
SET plan_bin = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdatePlanBinByID implements Querier.UpdatePlanBinByID.
func (q *DBQuerier) UpdatePlanBinByID(ctx context.Context, planBin []byte, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanBinByID")
	row := q.conn.QueryRow(ctx, updatePlanBinByIDSQL, planBin, runID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdatePlanBinByID: %w", err)
	}
	return item, nil
}

// UpdatePlanBinByIDBatch implements Querier.UpdatePlanBinByIDBatch.
func (q *DBQuerier) UpdatePlanBinByIDBatch(batch genericBatch, planBin []byte, runID pgtype.Text) {
	batch.Queue(updatePlanBinByIDSQL, planBin, runID)
}

// UpdatePlanBinByIDScan implements Querier.UpdatePlanBinByIDScan.
func (q *DBQuerier) UpdatePlanBinByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdatePlanBinByIDBatch row: %w", err)
	}
	return item, nil
}

const updatePlanJSONByIDSQL = `UPDATE plans
SET plan_json = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdatePlanJSONByID implements Querier.UpdatePlanJSONByID.
func (q *DBQuerier) UpdatePlanJSONByID(ctx context.Context, planJSON []byte, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdatePlanJSONByID")
	row := q.conn.QueryRow(ctx, updatePlanJSONByIDSQL, planJSON, runID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdatePlanJSONByID: %w", err)
	}
	return item, nil
}

// UpdatePlanJSONByIDBatch implements Querier.UpdatePlanJSONByIDBatch.
func (q *DBQuerier) UpdatePlanJSONByIDBatch(batch genericBatch, planJSON []byte, runID pgtype.Text) {
	batch.Queue(updatePlanJSONByIDSQL, planJSON, runID)
}

// UpdatePlanJSONByIDScan implements Querier.UpdatePlanJSONByIDScan.
func (q *DBQuerier) UpdatePlanJSONByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdatePlanJSONByIDBatch row: %w", err)
	}
	return item, nil
}
