package app

import (
	"context"
	"fmt"

	"github.com/go-logr/logr"
	"github.com/leg100/otf"
	"github.com/leg100/otf/inmem"
)

var _ otf.RunService = (*RunService)(nil)

type RunService struct {
	db    otf.DB
	cache otf.Cache
	proxy otf.ChunkStore
	otf.EventService
	logr.Logger
	*otf.RunFactory
}

func NewRunService(db otf.DB, logger logr.Logger, wss otf.WorkspaceService, cvs otf.ConfigurationVersionService, es otf.EventService, cache otf.Cache) (*RunService, error) {
	proxy, err := inmem.NewChunkProxy(cache, db)
	if err != nil {
		return nil, fmt.Errorf("constructing chunk proxy: %w", err)
	}
	return &RunService{
		db:           db,
		EventService: es,
		cache:        cache,
		proxy:        proxy,
		Logger:       logger,
		RunFactory: &otf.RunFactory{
			WorkspaceService:            wss,
			ConfigurationVersionService: cvs,
		},
	}, nil
}

// Create constructs and persists a new run object to the db, before scheduling
// the run.
func (s RunService) Create(ctx context.Context, spec otf.WorkspaceSpec, opts otf.RunCreateOptions) (*otf.Run, error) {
	run, err := s.New(ctx, spec, opts)
	if err != nil {
		s.Error(err, "constructing new run")
		return nil, err
	}

	if err = s.db.CreateRun(ctx, run); err != nil {
		s.Error(err, "creating run", "id", run.ID())
		return nil, err
	}
	s.V(1).Info("created run", "id", run.ID())

	s.Publish(otf.Event{Type: otf.EventRunCreated, Payload: run})

	return run, nil
}

// Get retrieves a run from the db.
func (s RunService) Get(ctx context.Context, runID string) (*otf.Run, error) {
	run, err := s.db.GetRun(ctx, runID)
	if err != nil {
		s.Error(err, "retrieving run", "id", runID)
		return nil, err
	}
	s.V(2).Info("retrieved run", "id", runID)

	return run, nil
}

// List retrieves multiple run objs. Use opts to filter and paginate the list.
func (s RunService) List(ctx context.Context, opts otf.RunListOptions) (*otf.RunList, error) {
	rl, err := s.db.ListRuns(ctx, opts)
	if err != nil {
		s.Error(err, "listing runs")
		return nil, err
	}

	s.V(2).Info("listed runs", append(opts.LogFields(), "count", len(rl.Items))...)

	return rl, nil
}

// ListWatch lists runs and then watches for changes to runs. Note: The options
// filter the list but not the watch.
func (s RunService) ListWatch(ctx context.Context, opts otf.RunListOptions) (<-chan *otf.Run, error) {
	// retrieve incomplete runs from db
	existing, err := s.db.ListRuns(ctx, opts)
	if err != nil {
		return nil, err
	}
	spool := make(chan *otf.Run, len(existing.Items))
	for _, r := range existing.Items {
		spool <- r
	}
	sub, err := s.Subscribe("run-listwatch")
	if err != nil {
		return nil, err
	}
	go func() {
		for {
			select {
			case <-ctx.Done():
				// context cancelled; shutdown spooler
				close(spool)
				return
			case event, ok := <-sub.C():
				if !ok {
					// sender closed channel; shutdown spooler
					close(spool)
					return
				}
				run, ok := event.Payload.(*otf.Run)
				if !ok {
					// skip non-run events
					continue
				}
				spool <- run
			}
		}
	}()
	return spool, nil
}

func (s RunService) Apply(ctx context.Context, runID string, opts otf.RunApplyOptions) error {
	run, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		return run.EnqueueApply()
	})
	if err != nil {
		s.Error(err, "enqueuing apply", "id", runID)
		return err
	}

	s.V(0).Info("enqueued apply", "id", runID)

	s.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})

	return err
}

func (s RunService) Discard(ctx context.Context, runID string, opts otf.RunDiscardOptions) error {
	run, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		return run.Discard()
	})
	if err != nil {
		s.Error(err, "discarding run", "id", runID)
		return err
	}

	s.V(0).Info("discarded run", "id", runID)

	s.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})

	return err
}

// Cancel a run. The run is canceled immediately if possible; otherwise a
// cancellation request is enqueued.
func (s RunService) Cancel(ctx context.Context, runID string, opts otf.RunCancelOptions) error {
	_, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		enqueue, err := run.Cancel()
		if err != nil {
			return err
		}
		if enqueue {
			s.V(0).Info("enqueued cancel request", "id", runID)
			// notify agent which'll send a SIGINT to terraform
			s.Publish(otf.Event{Type: otf.EventRunCancel, Payload: run})
		}
		s.V(0).Info("canceled run", "id", runID)
		s.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})
		return nil
	})
	if err != nil {
		s.Error(err, "canceling run", "id", runID)
		return err
	}
	return nil
}

// ForceCancel forcefully cancels a run.
func (s RunService) ForceCancel(ctx context.Context, runID string, opts otf.RunForceCancelOptions) error {
	run, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		return run.ForceCancel()
	})
	if err != nil {
		s.Error(err, "force canceling run", "id", runID)
		return err
	}
	s.V(0).Info("force canceled run", "id", runID)

	// notify agent which'll send a SIGKILL to terraform
	s.Publish(otf.Event{Type: otf.EventRunForceCancel, Payload: run})

	return err
}

func (s RunService) EnqueuePlan(ctx context.Context, runID string) (*otf.Run, error) {
	run, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		return run.EnqueuePlan()
	})
	if err != nil {
		s.Error(err, "started run", "id", runID)
		return nil, err
	}
	s.V(0).Info("started run", "id", runID)

	s.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})

	return run, nil
}

// GetPlanFile returns the plan file for the run.
func (s RunService) GetPlanFile(ctx context.Context, runID string, format otf.PlanFormat) ([]byte, error) {
	if plan, err := s.cache.Get(format.CacheKey(runID)); err == nil {
		return plan, nil
	}
	// Cache is empty; retrieve from DB
	file, err := s.db.GetPlanFile(ctx, runID, format)
	if err != nil {
		s.Error(err, "retrieving plan file", "id", runID, "format", format)
		return nil, err
	}
	// Cache plan before returning
	if err := s.cache.Set(format.CacheKey(runID), file); err != nil {
		return nil, fmt.Errorf("caching plan: %w", err)
	}
	return file, nil
}

// UploadPlanFile persists a run's plan file. The plan file is expected to have
// been produced using `terraform plan`. If the plan file is JSON serialized
// then its parsed for a summary of planned changes and the Plan object is
// updated accordingly.
func (s RunService) UploadPlanFile(ctx context.Context, runID string, plan []byte, format otf.PlanFormat) error {
	if err := s.db.SetPlanFile(ctx, runID, plan, format); err != nil {
		s.Error(err, "uploading plan file", "id", runID, "format", format)
		return err
	}

	s.V(0).Info("uploaded plan file", "id", runID, "format", format)

	if format == otf.PlanFormatJSON {
		report, err := otf.CompilePlanReport(plan)
		if err != nil {
			s.Error(err, "compiling planned changes report", "id", runID)
			return err
		}
		if err := s.db.CreatePlanReport(ctx, runID, report); err != nil {
			s.Error(err, "saving planned changes report", "id", runID)
			return err
		}
		s.V(1).Info("created planned changes report", "id", runID,
			"adds", report.Additions,
			"changes", report.Changes,
			"destructions", report.Destructions)
	}

	if err := s.cache.Set(format.CacheKey(runID), plan); err != nil {
		return fmt.Errorf("caching plan: %w", err)
	}

	return nil
}

// Delete deletes a terraform run.
func (s RunService) Delete(ctx context.Context, id string) error {
	// get run first so that we can include it in an event below
	run, err := s.db.GetRun(ctx, id)
	if err != nil {
		return err
	}
	if err := s.db.DeleteRun(ctx, id); err != nil {
		s.Error(err, "deleting run", "id", id)
		return err
	}
	s.V(0).Info("deleted run", "id", id)
	s.Publish(otf.Event{Type: otf.EventRunDeleted, Payload: run})
	return nil
}

// Start plan phase.
func (s RunService) Start(ctx context.Context, runID string, phase otf.PhaseType, opts otf.PhaseStartOptions) (*otf.Run, error) {
	run, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		return run.Start(phase)
	})
	if err != nil {
		s.Error(err, "starting phase", "id", runID, "phase", phase)
		return nil, err
	}
	s.V(0).Info("started phase", "id", runID, "phase", phase)
	s.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})
	return run, nil
}

// Finish plan phase.
func (s RunService) Finish(ctx context.Context, runID string, phase otf.PhaseType, opts otf.PhaseFinishOptions) (*otf.Run, error) {
	run, err := s.db.UpdateStatus(ctx, runID, func(run *otf.Run) error {
		return run.Finish(phase, opts)
	})
	if err != nil {
		s.Error(err, "finishing phase", "id", runID, "phase", phase)
		return nil, err
	}
	s.V(0).Info("finished phase", "id", runID, "phase", phase)
	s.Publish(otf.Event{Type: otf.EventRunStatusUpdate, Payload: run})
	return run, nil
}

// GetChunk reads a chunk of logs for a plan.
func (s RunService) GetChunk(ctx context.Context, runID string, phase otf.PhaseType, opts otf.GetChunkOptions) (otf.Chunk, error) {
	logs, err := s.proxy.GetChunk(ctx, runID, phase, opts)
	if err != nil {
		s.Error(err, "reading logs", "id", runID, "offset", opts.Offset, "limit", opts.Limit)
		return otf.Chunk{}, err
	}
	s.V(2).Info("read logs", "id", runID, "offset", opts.Offset, "limit", opts.Limit)
	return logs, nil
}

// PutChunk writes a chunk of logs for a plan.
func (s RunService) PutChunk(ctx context.Context, runID string, phase otf.PhaseType, chunk otf.Chunk) error {
	err := s.proxy.PutChunk(ctx, runID, phase, chunk)
	if err != nil {
		s.Error(err, "writing logs", "id", runID, "start", chunk.Start, "end", chunk.End)
		return err
	}
	s.V(2).Info("written logs", "id", runID, "start", chunk.Start, "end", chunk.End)
	return nil
}
