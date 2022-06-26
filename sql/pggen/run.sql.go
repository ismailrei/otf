// Code generated by pggen. DO NOT EDIT.

package pggen

import (
	"context"
	"fmt"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgtype"
	"github.com/jackc/pgx/v4"
)

const insertRunSQL = `INSERT INTO runs (
    run_id,
    created_at,
    is_destroy,
    position_in_queue,
    refresh,
    refresh_only,
    status,
    replace_addrs,
    target_addrs,
    configuration_version_id,
    workspace_id
) VALUES (
    $1,
    $2,
    $3,
    $4,
    $5,
    $6,
    $7,
    $8,
    $9,
    $10,
    $11
);`

type InsertRunParams struct {
	ID                     pgtype.Text
	CreatedAt              pgtype.Timestamptz
	IsDestroy              bool
	PositionInQueue        int
	Refresh                bool
	RefreshOnly            bool
	Status                 pgtype.Text
	ReplaceAddrs           []string
	TargetAddrs            []string
	ConfigurationVersionID pgtype.Text
	WorkspaceID            pgtype.Text
}

// InsertRun implements Querier.InsertRun.
func (q *DBQuerier) InsertRun(ctx context.Context, params InsertRunParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRun")
	cmdTag, err := q.conn.Exec(ctx, insertRunSQL, params.ID, params.CreatedAt, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Status, params.ReplaceAddrs, params.TargetAddrs, params.ConfigurationVersionID, params.WorkspaceID)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertRun: %w", err)
	}
	return cmdTag, err
}

// InsertRunBatch implements Querier.InsertRunBatch.
func (q *DBQuerier) InsertRunBatch(batch genericBatch, params InsertRunParams) {
	batch.Queue(insertRunSQL, params.ID, params.CreatedAt, params.IsDestroy, params.PositionInQueue, params.Refresh, params.RefreshOnly, params.Status, params.ReplaceAddrs, params.TargetAddrs, params.ConfigurationVersionID, params.WorkspaceID)
}

// InsertRunScan implements Querier.InsertRunScan.
func (q *DBQuerier) InsertRunScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertRunBatch: %w", err)
	}
	return cmdTag, err
}

const insertRunStatusTimestampSQL = `INSERT INTO run_status_timestamps (
    run_id,
    status,
    timestamp
) VALUES (
    $1,
    $2,
    $3
);`

type InsertRunStatusTimestampParams struct {
	ID        pgtype.Text
	Status    pgtype.Text
	Timestamp pgtype.Timestamptz
}

// InsertRunStatusTimestamp implements Querier.InsertRunStatusTimestamp.
func (q *DBQuerier) InsertRunStatusTimestamp(ctx context.Context, params InsertRunStatusTimestampParams) (pgconn.CommandTag, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "InsertRunStatusTimestamp")
	cmdTag, err := q.conn.Exec(ctx, insertRunStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
	if err != nil {
		return cmdTag, fmt.Errorf("exec query InsertRunStatusTimestamp: %w", err)
	}
	return cmdTag, err
}

// InsertRunStatusTimestampBatch implements Querier.InsertRunStatusTimestampBatch.
func (q *DBQuerier) InsertRunStatusTimestampBatch(batch genericBatch, params InsertRunStatusTimestampParams) {
	batch.Queue(insertRunStatusTimestampSQL, params.ID, params.Status, params.Timestamp)
}

// InsertRunStatusTimestampScan implements Querier.InsertRunStatusTimestampScan.
func (q *DBQuerier) InsertRunStatusTimestampScan(results pgx.BatchResults) (pgconn.CommandTag, error) {
	cmdTag, err := results.Exec()
	if err != nil {
		return cmdTag, fmt.Errorf("exec InsertRunStatusTimestampBatch: %w", err)
	}
	return cmdTag, err
}

const findRunsSQL = `SELECT
    runs.run_id,
    runs.created_at,
    runs.force_cancel_available_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.status,
    plans.status      AS plan_status,
    applies.status      AS apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    plans.report AS planned_changes,
    applies.report AS applied_changes,
    runs.configuration_version_id,
    runs.workspace_id,
    configuration_versions.speculative,
    workspaces.auto_apply,
    workspaces.name AS workspace_name,
    organizations.name AS organization_name,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = plans.run_id
        AND   st.phase = 'plan'
        GROUP BY run_id, phase
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = applies.run_id
        AND   st.phase = 'apply'
        GROUP BY run_id, phase
    ) AS apply_status_timestamps
FROM runs
JOIN plans USING (run_id)
JOIN applies USING (run_id)
JOIN configuration_versions USING(configuration_version_id)
JOIN workspaces ON runs.workspace_id = workspaces.workspace_id
JOIN organizations USING(organization_id)
WHERE
    organizations.name      LIKE ANY($1)
AND workspaces.workspace_id LIKE ANY($2)
AND workspaces.name         LIKE ANY($3)
AND runs.status             LIKE ANY($4)
ORDER BY runs.created_at ASC
LIMIT $5 OFFSET $6
;`

type FindRunsParams struct {
	OrganizationNames []string
	WorkspaceIds      []string
	WorkspaceNames    []string
	Statuses          []string
	Limit             int
	Offset            int
}

type FindRunsRow struct {
	RunID                  pgtype.Text             `json:"run_id"`
	CreatedAt              pgtype.Timestamptz      `json:"created_at"`
	ForceCancelAvailableAt pgtype.Timestamptz      `json:"force_cancel_available_at"`
	IsDestroy              bool                    `json:"is_destroy"`
	PositionInQueue        int                     `json:"position_in_queue"`
	Refresh                bool                    `json:"refresh"`
	RefreshOnly            bool                    `json:"refresh_only"`
	Status                 pgtype.Text             `json:"status"`
	PlanStatus             pgtype.Text             `json:"plan_status"`
	ApplyStatus            pgtype.Text             `json:"apply_status"`
	ReplaceAddrs           []string                `json:"replace_addrs"`
	TargetAddrs            []string                `json:"target_addrs"`
	PlannedChanges         *Report                 `json:"planned_changes"`
	AppliedChanges         *Report                 `json:"applied_changes"`
	ConfigurationVersionID pgtype.Text             `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text             `json:"workspace_id"`
	Speculative            bool                    `json:"speculative"`
	AutoApply              bool                    `json:"auto_apply"`
	WorkspaceName          pgtype.Text             `json:"workspace_name"`
	OrganizationName       pgtype.Text             `json:"organization_name"`
	RunStatusTimestamps    []RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []PhaseStatusTimestamps `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []PhaseStatusTimestamps `json:"apply_status_timestamps"`
}

// FindRuns implements Querier.FindRuns.
func (q *DBQuerier) FindRuns(ctx context.Context, params FindRunsParams) ([]FindRunsRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRuns")
	rows, err := q.conn.Query(ctx, findRunsSQL, params.OrganizationNames, params.WorkspaceIds, params.WorkspaceNames, params.Statuses, params.Limit, params.Offset)
	if err != nil {
		return nil, fmt.Errorf("query FindRuns: %w", err)
	}
	defer rows.Close()
	items := []FindRunsRow{}
	plannedChangesRow := q.types.newReport()
	appliedChangesRow := q.types.newReport()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	for rows.Next() {
		var item FindRunsRow
		if err := rows.Scan(&item.RunID, &item.CreatedAt, &item.ForceCancelAvailableAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, plannedChangesRow, appliedChangesRow, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, &item.WorkspaceName, &item.OrganizationName, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
			return nil, fmt.Errorf("scan FindRuns row: %w", err)
		}
		if err := plannedChangesRow.AssignTo(&item.PlannedChanges); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := appliedChangesRow.AssignTo(&item.AppliedChanges); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRuns rows: %w", err)
	}
	return items, err
}

// FindRunsBatch implements Querier.FindRunsBatch.
func (q *DBQuerier) FindRunsBatch(batch genericBatch, params FindRunsParams) {
	batch.Queue(findRunsSQL, params.OrganizationNames, params.WorkspaceIds, params.WorkspaceNames, params.Statuses, params.Limit, params.Offset)
}

// FindRunsScan implements Querier.FindRunsScan.
func (q *DBQuerier) FindRunsScan(results pgx.BatchResults) ([]FindRunsRow, error) {
	rows, err := results.Query()
	if err != nil {
		return nil, fmt.Errorf("query FindRunsBatch: %w", err)
	}
	defer rows.Close()
	items := []FindRunsRow{}
	plannedChangesRow := q.types.newReport()
	appliedChangesRow := q.types.newReport()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	for rows.Next() {
		var item FindRunsRow
		if err := rows.Scan(&item.RunID, &item.CreatedAt, &item.ForceCancelAvailableAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, plannedChangesRow, appliedChangesRow, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, &item.WorkspaceName, &item.OrganizationName, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
			return nil, fmt.Errorf("scan FindRunsBatch row: %w", err)
		}
		if err := plannedChangesRow.AssignTo(&item.PlannedChanges); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := appliedChangesRow.AssignTo(&item.AppliedChanges); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
			return nil, fmt.Errorf("assign FindRuns row: %w", err)
		}
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("close FindRunsBatch rows: %w", err)
	}
	return items, err
}

const countRunsSQL = `SELECT count(*)
FROM runs
JOIN workspaces USING(workspace_id)
JOIN organizations USING(organization_id)
WHERE
    organizations.name      LIKE ANY($1)
AND workspaces.workspace_id LIKE ANY($2)
AND workspaces.name         LIKE ANY($3)
AND runs.status             LIKE ANY($4)
;`

type CountRunsParams struct {
	OrganizationNames []string
	WorkspaceIds      []string
	WorkspaceNames    []string
	Statuses          []string
}

// CountRuns implements Querier.CountRuns.
func (q *DBQuerier) CountRuns(ctx context.Context, params CountRunsParams) (*int, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "CountRuns")
	row := q.conn.QueryRow(ctx, countRunsSQL, params.OrganizationNames, params.WorkspaceIds, params.WorkspaceNames, params.Statuses)
	var item int
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("query CountRuns: %w", err)
	}
	return &item, nil
}

// CountRunsBatch implements Querier.CountRunsBatch.
func (q *DBQuerier) CountRunsBatch(batch genericBatch, params CountRunsParams) {
	batch.Queue(countRunsSQL, params.OrganizationNames, params.WorkspaceIds, params.WorkspaceNames, params.Statuses)
}

// CountRunsScan implements Querier.CountRunsScan.
func (q *DBQuerier) CountRunsScan(results pgx.BatchResults) (*int, error) {
	row := results.QueryRow()
	var item int
	if err := row.Scan(&item); err != nil {
		return &item, fmt.Errorf("scan CountRunsBatch row: %w", err)
	}
	return &item, nil
}

const findRunByIDSQL = `SELECT
    runs.run_id,
    runs.created_at,
    runs.force_cancel_available_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.status,
    plans.status      AS plan_status,
    applies.status      AS apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    plans.report AS planned_changes,
    applies.report AS applied_changes,
    runs.configuration_version_id,
    runs.workspace_id,
    configuration_versions.speculative,
    workspaces.auto_apply,
    workspaces.name AS workspace_name,
    organizations.name AS organization_name,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = plans.run_id
        AND   st.phase = 'plan'
        GROUP BY run_id, phase
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = applies.run_id
        AND   st.phase = 'apply'
        GROUP BY run_id, phase
    ) AS apply_status_timestamps
FROM runs
JOIN plans USING (run_id)
JOIN applies USING (run_id)
JOIN configuration_versions USING(configuration_version_id)
JOIN workspaces ON runs.workspace_id = workspaces.workspace_id
JOIN organizations USING(organization_id)
WHERE runs.run_id = $1
;`

type FindRunByIDRow struct {
	RunID                  pgtype.Text             `json:"run_id"`
	CreatedAt              pgtype.Timestamptz      `json:"created_at"`
	ForceCancelAvailableAt pgtype.Timestamptz      `json:"force_cancel_available_at"`
	IsDestroy              bool                    `json:"is_destroy"`
	PositionInQueue        int                     `json:"position_in_queue"`
	Refresh                bool                    `json:"refresh"`
	RefreshOnly            bool                    `json:"refresh_only"`
	Status                 pgtype.Text             `json:"status"`
	PlanStatus             pgtype.Text             `json:"plan_status"`
	ApplyStatus            pgtype.Text             `json:"apply_status"`
	ReplaceAddrs           []string                `json:"replace_addrs"`
	TargetAddrs            []string                `json:"target_addrs"`
	PlannedChanges         *Report                 `json:"planned_changes"`
	AppliedChanges         *Report                 `json:"applied_changes"`
	ConfigurationVersionID pgtype.Text             `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text             `json:"workspace_id"`
	Speculative            bool                    `json:"speculative"`
	AutoApply              bool                    `json:"auto_apply"`
	WorkspaceName          pgtype.Text             `json:"workspace_name"`
	OrganizationName       pgtype.Text             `json:"organization_name"`
	RunStatusTimestamps    []RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []PhaseStatusTimestamps `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []PhaseStatusTimestamps `json:"apply_status_timestamps"`
}

// FindRunByID implements Querier.FindRunByID.
func (q *DBQuerier) FindRunByID(ctx context.Context, runID pgtype.Text) (FindRunByIDRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByID")
	row := q.conn.QueryRow(ctx, findRunByIDSQL, runID)
	var item FindRunByIDRow
	plannedChangesRow := q.types.newReport()
	appliedChangesRow := q.types.newReport()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.ForceCancelAvailableAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, plannedChangesRow, appliedChangesRow, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, &item.WorkspaceName, &item.OrganizationName, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("query FindRunByID: %w", err)
	}
	if err := plannedChangesRow.AssignTo(&item.PlannedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := appliedChangesRow.AssignTo(&item.AppliedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	return item, nil
}

// FindRunByIDBatch implements Querier.FindRunByIDBatch.
func (q *DBQuerier) FindRunByIDBatch(batch genericBatch, runID pgtype.Text) {
	batch.Queue(findRunByIDSQL, runID)
}

// FindRunByIDScan implements Querier.FindRunByIDScan.
func (q *DBQuerier) FindRunByIDScan(results pgx.BatchResults) (FindRunByIDRow, error) {
	row := results.QueryRow()
	var item FindRunByIDRow
	plannedChangesRow := q.types.newReport()
	appliedChangesRow := q.types.newReport()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.ForceCancelAvailableAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, plannedChangesRow, appliedChangesRow, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, &item.WorkspaceName, &item.OrganizationName, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("scan FindRunByIDBatch row: %w", err)
	}
	if err := plannedChangesRow.AssignTo(&item.PlannedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := appliedChangesRow.AssignTo(&item.AppliedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByID row: %w", err)
	}
	return item, nil
}

const findRunByIDForUpdateSQL = `SELECT
    runs.run_id,
    runs.created_at,
    runs.force_cancel_available_at,
    runs.is_destroy,
    runs.position_in_queue,
    runs.refresh,
    runs.refresh_only,
    runs.status,
    plans.status        AS plan_status,
    applies.status      AS apply_status,
    runs.replace_addrs,
    runs.target_addrs,
    plans.report AS planned_changes,
    applies.report AS applied_changes,
    runs.configuration_version_id,
    runs.workspace_id,
    configuration_versions.speculative,
    workspaces.auto_apply,
    workspaces.name AS workspace_name,
    organizations.name AS organization_name,
    (
        SELECT array_agg(rst.*) AS run_status_timestamps
        FROM run_status_timestamps rst
        WHERE rst.run_id = runs.run_id
        GROUP BY run_id
    ) AS run_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = plans.run_id
        AND   st.phase = 'plan'
        GROUP BY run_id, phase
    ) AS plan_status_timestamps,
    (
        SELECT array_agg(st.*) AS phase_status_timestamps
        FROM phase_status_timestamps st
        WHERE st.run_id = applies.run_id
        AND   st.phase = 'apply'
        GROUP BY run_id, phase
    ) AS apply_status_timestamps
FROM runs
JOIN plans USING (run_id)
JOIN applies USING (run_id)
JOIN configuration_versions USING(configuration_version_id)
JOIN workspaces ON runs.workspace_id = workspaces.workspace_id
JOIN organizations USING(organization_id)
WHERE runs.run_id = $1
FOR UPDATE
;`

type FindRunByIDForUpdateRow struct {
	RunID                  pgtype.Text             `json:"run_id"`
	CreatedAt              pgtype.Timestamptz      `json:"created_at"`
	ForceCancelAvailableAt pgtype.Timestamptz      `json:"force_cancel_available_at"`
	IsDestroy              bool                    `json:"is_destroy"`
	PositionInQueue        int                     `json:"position_in_queue"`
	Refresh                bool                    `json:"refresh"`
	RefreshOnly            bool                    `json:"refresh_only"`
	Status                 pgtype.Text             `json:"status"`
	PlanStatus             pgtype.Text             `json:"plan_status"`
	ApplyStatus            pgtype.Text             `json:"apply_status"`
	ReplaceAddrs           []string                `json:"replace_addrs"`
	TargetAddrs            []string                `json:"target_addrs"`
	PlannedChanges         *Report                 `json:"planned_changes"`
	AppliedChanges         *Report                 `json:"applied_changes"`
	ConfigurationVersionID pgtype.Text             `json:"configuration_version_id"`
	WorkspaceID            pgtype.Text             `json:"workspace_id"`
	Speculative            bool                    `json:"speculative"`
	AutoApply              bool                    `json:"auto_apply"`
	WorkspaceName          pgtype.Text             `json:"workspace_name"`
	OrganizationName       pgtype.Text             `json:"organization_name"`
	RunStatusTimestamps    []RunStatusTimestamps   `json:"run_status_timestamps"`
	PlanStatusTimestamps   []PhaseStatusTimestamps `json:"plan_status_timestamps"`
	ApplyStatusTimestamps  []PhaseStatusTimestamps `json:"apply_status_timestamps"`
}

// FindRunByIDForUpdate implements Querier.FindRunByIDForUpdate.
func (q *DBQuerier) FindRunByIDForUpdate(ctx context.Context, runID pgtype.Text) (FindRunByIDForUpdateRow, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "FindRunByIDForUpdate")
	row := q.conn.QueryRow(ctx, findRunByIDForUpdateSQL, runID)
	var item FindRunByIDForUpdateRow
	plannedChangesRow := q.types.newReport()
	appliedChangesRow := q.types.newReport()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.ForceCancelAvailableAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, plannedChangesRow, appliedChangesRow, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, &item.WorkspaceName, &item.OrganizationName, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("query FindRunByIDForUpdate: %w", err)
	}
	if err := plannedChangesRow.AssignTo(&item.PlannedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := appliedChangesRow.AssignTo(&item.AppliedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	return item, nil
}

// FindRunByIDForUpdateBatch implements Querier.FindRunByIDForUpdateBatch.
func (q *DBQuerier) FindRunByIDForUpdateBatch(batch genericBatch, runID pgtype.Text) {
	batch.Queue(findRunByIDForUpdateSQL, runID)
}

// FindRunByIDForUpdateScan implements Querier.FindRunByIDForUpdateScan.
func (q *DBQuerier) FindRunByIDForUpdateScan(results pgx.BatchResults) (FindRunByIDForUpdateRow, error) {
	row := results.QueryRow()
	var item FindRunByIDForUpdateRow
	plannedChangesRow := q.types.newReport()
	appliedChangesRow := q.types.newReport()
	runStatusTimestampsArray := q.types.newRunStatusTimestampsArray()
	planStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	applyStatusTimestampsArray := q.types.newPhaseStatusTimestampsArray()
	if err := row.Scan(&item.RunID, &item.CreatedAt, &item.ForceCancelAvailableAt, &item.IsDestroy, &item.PositionInQueue, &item.Refresh, &item.RefreshOnly, &item.Status, &item.PlanStatus, &item.ApplyStatus, &item.ReplaceAddrs, &item.TargetAddrs, plannedChangesRow, appliedChangesRow, &item.ConfigurationVersionID, &item.WorkspaceID, &item.Speculative, &item.AutoApply, &item.WorkspaceName, &item.OrganizationName, runStatusTimestampsArray, planStatusTimestampsArray, applyStatusTimestampsArray); err != nil {
		return item, fmt.Errorf("scan FindRunByIDForUpdateBatch row: %w", err)
	}
	if err := plannedChangesRow.AssignTo(&item.PlannedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := appliedChangesRow.AssignTo(&item.AppliedChanges); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := runStatusTimestampsArray.AssignTo(&item.RunStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := planStatusTimestampsArray.AssignTo(&item.PlanStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	if err := applyStatusTimestampsArray.AssignTo(&item.ApplyStatusTimestamps); err != nil {
		return item, fmt.Errorf("assign FindRunByIDForUpdate row: %w", err)
	}
	return item, nil
}

const updateRunStatusSQL = `UPDATE runs
SET
    status = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateRunStatus implements Querier.UpdateRunStatus.
func (q *DBQuerier) UpdateRunStatus(ctx context.Context, status pgtype.Text, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunStatus")
	row := q.conn.QueryRow(ctx, updateRunStatusSQL, status, id)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateRunStatus: %w", err)
	}
	return item, nil
}

// UpdateRunStatusBatch implements Querier.UpdateRunStatusBatch.
func (q *DBQuerier) UpdateRunStatusBatch(batch genericBatch, status pgtype.Text, id pgtype.Text) {
	batch.Queue(updateRunStatusSQL, status, id)
}

// UpdateRunStatusScan implements Querier.UpdateRunStatusScan.
func (q *DBQuerier) UpdateRunStatusScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateRunStatusBatch row: %w", err)
	}
	return item, nil
}

const updateRunForceCancelAvailableAtSQL = `UPDATE runs
SET
    force_cancel_available_at = $1
WHERE run_id = $2
RETURNING run_id
;`

// UpdateRunForceCancelAvailableAt implements Querier.UpdateRunForceCancelAvailableAt.
func (q *DBQuerier) UpdateRunForceCancelAvailableAt(ctx context.Context, forceCancelAvailableAt pgtype.Timestamptz, id pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "UpdateRunForceCancelAvailableAt")
	row := q.conn.QueryRow(ctx, updateRunForceCancelAvailableAtSQL, forceCancelAvailableAt, id)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query UpdateRunForceCancelAvailableAt: %w", err)
	}
	return item, nil
}

// UpdateRunForceCancelAvailableAtBatch implements Querier.UpdateRunForceCancelAvailableAtBatch.
func (q *DBQuerier) UpdateRunForceCancelAvailableAtBatch(batch genericBatch, forceCancelAvailableAt pgtype.Timestamptz, id pgtype.Text) {
	batch.Queue(updateRunForceCancelAvailableAtSQL, forceCancelAvailableAt, id)
}

// UpdateRunForceCancelAvailableAtScan implements Querier.UpdateRunForceCancelAvailableAtScan.
func (q *DBQuerier) UpdateRunForceCancelAvailableAtScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan UpdateRunForceCancelAvailableAtBatch row: %w", err)
	}
	return item, nil
}

const deleteRunByIDSQL = `DELETE
FROM runs
WHERE run_id = $1
RETURNING run_id
;`

// DeleteRunByID implements Querier.DeleteRunByID.
func (q *DBQuerier) DeleteRunByID(ctx context.Context, runID pgtype.Text) (pgtype.Text, error) {
	ctx = context.WithValue(ctx, "pggen_query_name", "DeleteRunByID")
	row := q.conn.QueryRow(ctx, deleteRunByIDSQL, runID)
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("query DeleteRunByID: %w", err)
	}
	return item, nil
}

// DeleteRunByIDBatch implements Querier.DeleteRunByIDBatch.
func (q *DBQuerier) DeleteRunByIDBatch(batch genericBatch, runID pgtype.Text) {
	batch.Queue(deleteRunByIDSQL, runID)
}

// DeleteRunByIDScan implements Querier.DeleteRunByIDScan.
func (q *DBQuerier) DeleteRunByIDScan(results pgx.BatchResults) (pgtype.Text, error) {
	row := results.QueryRow()
	var item pgtype.Text
	if err := row.Scan(&item); err != nil {
		return item, fmt.Errorf("scan DeleteRunByIDBatch row: %w", err)
	}
	return item, nil
}
