package site

import (
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

func TestNewAppRoutes(t *testing.T) {
	prevFetcher := githubChartFetcher
	prevChart := defaultSite.chart
	prevChartAt := defaultSite.chartAt
	prevChartRetry := defaultSite.chartRetry

	githubChartFetcher = func(string, time.Time) (*GitHubChart, error) {
		months, weeks, rows := buildGitHubChartLayout([]GitHubChartDay{
			{Date: "2026-01-01", Count: 3, Level: 1, Tooltip: "3 contributions on January 1st."},
			{Date: "2026-01-02", Count: 8, Level: 2, Tooltip: "8 contributions on January 2nd."},
			{Date: "2026-01-03", Count: 12, Level: 3, Tooltip: "12 contributions on January 3rd."},
			{Date: "2026-01-04", Count: 19, Level: 4, Tooltip: "19 contributions on January 4th."},
			{Date: "2026-01-05", Count: 0, Level: 0, Tooltip: "No contributions on January 5th."},
		}, 2026)

		return &GitHubChart{
			Username:   "0mjs",
			ProfileURL: "https://github.com/0mjs",
			Year:       2026,
			Total:      42,
			Months:     months,
			Weeks:      weeks,
			Rows:       rows,
		}, nil
	}
	defaultSite.chart = nil
	defaultSite.chartAt = time.Time{}
	defaultSite.chartRetry = time.Time{}
	t.Cleanup(func() {
		githubChartFetcher = prevFetcher
		defaultSite.chart = prevChart
		defaultSite.chartAt = prevChartAt
		defaultSite.chartRetry = prevChartRetry
	})

	app := NewApp()

	tests := []struct {
		name       string
		path       string
		statusCode int
		bodyHas    string
		headers    map[string]string
	}{
		{
			name:       "home",
			path:       "/",
			statusCode: http.StatusOK,
			bodyHas:    "42 contributions in 2026",
		},
		{
			name:       "about",
			path:       "/about",
			statusCode: http.StatusOK,
			bodyHas:    "Hi, I&rsquo;m Matt",
		},
		{
			name:       "post",
			path:       "/post/why-go",
			statusCode: http.StatusOK,
			bodyHas:    "Why Go?",
		},
		{
			name:       "missing post",
			path:       "/post/missing",
			statusCode: http.StatusNotFound,
			bodyHas:    "Not Found",
		},
		{
			name:       "static css",
			path:       "/css/style.css",
			statusCode: http.StatusOK,
			bodyHas:    ".container",
			headers: map[string]string{
				"Cache-Control": "public, max-age=604800",
				"Content-Type":  "text/css; charset=utf-8",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rec := httptest.NewRecorder()

			app.ServeHTTP(rec, req)

			if rec.Code != tt.statusCode {
				t.Fatalf("status=%d want=%d body=%q", rec.Code, tt.statusCode, rec.Body.String())
			}

			body := rec.Body.String()
			if tt.bodyHas != "" && !strings.Contains(body, tt.bodyHas) {
				t.Fatalf("body does not contain %q: %q", tt.bodyHas, body)
			}

			for key, want := range tt.headers {
				if got := rec.Header().Get(key); got != want {
					t.Fatalf("%s=%q want=%q", key, got, want)
				}
			}
		})
	}
}

func TestBuildGitHubChartLayoutAlignsWeeks(t *testing.T) {
	months, weeks, rows := buildGitHubChartLayout([]GitHubChartDay{
		{Date: "2026-01-01", Count: 5, Level: 1, Tooltip: "5 contributions on January 1st."},
		{Date: "2026-01-02", Count: 4, Level: 1, Tooltip: "4 contributions on January 2nd."},
		{Date: "2026-01-03", Count: 3, Level: 1, Tooltip: "3 contributions on January 3rd."},
		{Date: "2026-01-04", Count: 2, Level: 1, Tooltip: "2 contributions on January 4th."},
	}, 2026)

	if len(weeks) != 53 {
		t.Fatalf("weeks=%d", len(weeks))
	}

	for dayIndex := 0; dayIndex < 4; dayIndex++ {
		if !weeks[0].Days[dayIndex].Empty {
			t.Fatalf("first week day %d should be empty", dayIndex)
		}
	}

	if got := weeks[0].Days[4].Date; got != "2026-01-01" {
		t.Fatalf("first week thursday=%q", got)
	}
	if got := weeks[0].Days[5].Date; got != "2026-01-02" {
		t.Fatalf("first week friday=%q", got)
	}
	if got := weeks[0].Days[6].Date; got != "2026-01-03" {
		t.Fatalf("first week saturday=%q", got)
	}
	if got := weeks[1].Days[0].Date; got != "2026-01-04" {
		t.Fatalf("second week sunday=%q", got)
	}

	if len(months) < 2 {
		t.Fatalf("months=%d", len(months))
	}
	if months[0].Name != "Jan" || months[0].ColSpan != 5 {
		t.Fatalf("first month=%+v", months[0])
	}
	if months[1].Name != "Feb" || months[1].ColSpan != 4 {
		t.Fatalf("second month=%+v", months[1])
	}

	if len(rows) != 7 {
		t.Fatalf("rows=%d", len(rows))
	}
	if rows[1].Label != "Mon" || rows[3].Label != "Wed" || rows[5].Label != "Fri" {
		t.Fatalf("row labels=%q,%q,%q", rows[1].Label, rows[3].Label, rows[5].Label)
	}
	if got := rows[4].Days[0].Date; got != "2026-01-01" {
		t.Fatalf("thursday row first week=%q", got)
	}
}
