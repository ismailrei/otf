package repo

import (
	"context"
	"net/http"
	"testing"

	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/leg100/otf/internal"
	"github.com/leg100/otf/internal/cloud"
	"github.com/leg100/otf/internal/sql"
	"github.com/stretchr/testify/require"
)

func newTestHook(t *testing.T, f factory, cloudID *string) *hook {
	want, err := f.newHook(newHookOpts{
		id:         internal.UUID(uuid.New()),
		secret:     internal.String("top-secret"),
		identifier: "leg100/" + uuid.NewString(),
		cloud:      "github",
		cloudID:    cloudID,
	})
	require.NoError(t, err)
	return want
}

func newTestFactory(t *testing.T, event cloud.VCSEvent) factory {
	return newFactory(
		fakeHostnameService{},
		fakeCloudService{event: event},
	)
}

func newTestDB(t *testing.T) *pgdb {
	db, err := sql.New(context.Background(), sql.Options{
		Logger:     logr.Discard(),
		ConnString: sql.NewTestDB(t),
	})
	require.NoError(t, err)
	t.Cleanup(db.Close)

	return &pgdb{
		DB: db,
		factory: factory{
			Service:         fakeCloudService{},
			HostnameService: fakeHostnameService{},
		},
	}
}

type fakeCloudService struct {
	event cloud.VCSEvent
	cloud.Service
}

func (f fakeCloudService) GetCloudConfig(string) (cloud.Config, error) {
	return cloud.Config{Cloud: &fakeCloud{event: f.event}}, nil
}

type fakeCloud struct {
	event cloud.VCSEvent

	cloud.Cloud
}

func (f *fakeCloud) HandleEvent(http.ResponseWriter, *http.Request, cloud.HandleEventOptions) cloud.VCSEvent {
	return f.event
}

type fakeHostnameService struct {
	hostname string

	internal.HostnameService
}

func (f fakeHostnameService) Hostname() string { return f.hostname }

type fakeCloudClient struct {
	hook      cloud.Webhook // seed cloud with hook
	gotUpdate bool

	cloud.Client
}

func (f *fakeCloudClient) CreateWebhook(context.Context, cloud.CreateWebhookOptions) (string, error) {
	return f.hook.ID, nil
}

func (f *fakeCloudClient) GetWebhook(ctx context.Context, opts cloud.GetWebhookOptions) (cloud.Webhook, error) {
	if f.hook.ID == opts.ID {
		return f.hook, nil
	}
	return cloud.Webhook{}, internal.ErrResourceNotFound
}

func (f *fakeCloudClient) UpdateWebhook(context.Context, cloud.UpdateWebhookOptions) error {
	f.gotUpdate = true

	return nil
}