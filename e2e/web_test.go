package e2e

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"strconv"
	"strings"
	"testing"

	"github.com/chromedp/cdproto/input"
	"github.com/chromedp/chromedp"
	"github.com/leg100/otf"
	"github.com/leg100/otf/http/html"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWeb(t *testing.T) {
	headless := true
	if v, ok := os.LookupEnv("OTF_E2E_HEADLESS"); ok {
		var err error
		headless, err = strconv.ParseBool(v)
		require.NoError(t, err)
	}

	githubHostname := githubStub(t)
	t.Setenv("OTF_GITHUB_HOSTNAME", githubHostname)

	url := startDaemon(t)

	allocCtx, cancel := chromedp.NewExecAllocator(context.Background(),
		append(chromedp.DefaultExecAllocatorOptions[:],
			chromedp.Flag("headless", headless),
			chromedp.Flag("hide-scrollbars", true),
			chromedp.Flag("mute-audio", true),
			chromedp.Flag("ignore-certificate-errors", true),
			chromedp.Flag("disable-gpu", true),
		)...)
	defer cancel()

	t.Run("login", func(t *testing.T) {
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		var gotLoginPrompt string
		var gotLocationOrganizations string

		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(url),
			screenshot("otf_login"),
			chromedp.Text(".center", &gotLoginPrompt, chromedp.NodeVisible),
			chromedp.Click(".login-button-github", chromedp.NodeVisible),
			screenshot("otf_login_successful"),
			chromedp.Location(&gotLocationOrganizations),
		})
		require.NoError(t, err)

		assert.Equal(t, "Login with Github", strings.TrimSpace(gotLoginPrompt))
		assert.Equal(t, url+"/organizations", gotLocationOrganizations)
	})

	t.Run("new workspace", func(t *testing.T) {
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		var gotFlashSuccess string
		workspaceName := "workspace-" + otf.GenerateRandomString(4)

		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(url),
			chromedp.Click(".login-button-github", chromedp.NodeVisible),
			chromedp.Click(".content-list a", chromedp.NodeVisible),
			chromedp.Click("#workspaces > a", chromedp.NodeVisible),
			chromedp.Click("#new-workspace-button", chromedp.NodeVisible),
			screenshot("otf_new_workspace_form"),
			chromedp.Focus("input#name", chromedp.NodeVisible),
			input.InsertText(workspaceName),
			chromedp.Click("#create-workspace-button"),
			screenshot("otf_created_workspace"),
			chromedp.Text(".flash-success", &gotFlashSuccess, chromedp.NodeVisible),
		})
		require.NoError(t, err)

		assert.Equal(t, "created workspace: "+workspaceName, strings.TrimSpace(gotFlashSuccess))
	})

	t.Run("list users", func(t *testing.T) {
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		var gotUser string
		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(url),
			// login
			chromedp.Click(".login-button-github", chromedp.NodeVisible),
			// select fake-user's personal org
			chromedp.Click("#item-organization-fake-user a", chromedp.NodeVisible),
			// list users
			chromedp.Click("#users > a", chromedp.NodeVisible),
			screenshot("otf_list_users"),
			chromedp.Text("#item-user-fake-user .status", &gotUser, chromedp.NodeVisible),
		})
		require.NoError(t, err)

		assert.Equal(t, "fake-user", strings.TrimSpace(gotUser))
	})

	t.Run("list team members", func(t *testing.T) {
		ctx, cancel := chromedp.NewContext(allocCtx)
		defer cancel()

		var gotUser string
		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.Navigate(url),
			// login
			chromedp.Click(".login-button-github", chromedp.NodeVisible),
			// select fake-user's personal org
			chromedp.Click("#item-organization-fake-user a", chromedp.NodeVisible),
			// list teams
			chromedp.Click("#teams > a", chromedp.NodeVisible),
			screenshot("otf_list_teams"),
			// select owners team
			chromedp.Click("#item-team-owners a", chromedp.NodeVisible),
			screenshot("otf_list_team_members"),
			chromedp.Text("#item-user-fake-user .status", &gotUser, chromedp.NodeVisible),
		})
		require.NoError(t, err)

		assert.Equal(t, "fake-user", strings.TrimSpace(gotUser))
	})
}

var screenshotCounter = 0

func screenshot(name string) chromedp.ActionFunc {
	return func(ctx context.Context) error {
		screenshotCounter++

		var image []byte
		err := chromedp.Run(ctx, chromedp.Tasks{
			chromedp.WaitReady(`body`),
			chromedp.CaptureScreenshot(&image),
		})
		if err != nil {
			return err
		}
		err = os.MkdirAll("screenshots", 0o755)
		if err != nil {
			return err
		}
		err = os.WriteFile(fmt.Sprintf("screenshots/%02d_%s.png", screenshotCounter, name), image, 0o644)
		if err != nil {
			return err
		}
		return nil
	}
}

func githubStub(t *testing.T) string {
	org := otf.NewTestOrganization(t)
	team := otf.NewTeam("fake-team", org)
	user := otf.NewUser("fake-user", otf.WithOrganizationMemberships(org), otf.WithTeamMemberships(team))
	srv := html.NewTestGithubServer(t, user)

	u, err := url.Parse(srv.URL)
	require.NoError(t, err)
	return u.Host
}