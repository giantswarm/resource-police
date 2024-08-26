package report

import (
	"embed"
	"errors"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/giantswarm/resource-police/pkg/cortex"
)

//go:embed golden/*
var fs embed.FS

func Test_RenderReport(t *testing.T) {
	tests := []struct {
		name       string
		clusters   []cortex.Cluster
		errors     []error
		goldenFile string
	}{
		{
			name: "Success",
			clusters: []cortex.Cluster{
				// 2 hours old cluster
				{
					Installation:         "foobar",
					ID:                   "r72ux",
					Release:              "14.5.6",
					Provider:             "openstack",
					FirstTimestamp:       time.Now().UTC().Add(-2 * time.Hour),
					NamespaceDescription: "mynamespace1",
				},
				// 4 hours old cluster
				{
					Installation:         "gaia",
					ID:                   "def34",
					Release:              "1.2.3",
					Provider:             "aws",
					FirstTimestamp:       time.Now().UTC().Add(-4 * time.Hour),
					NamespaceDescription: "mynamespace2",
				},
				// 2 days old cluster
				{
					Installation:         "blah",
					ID:                   "t2ayx",
					Release:              "13.0.4",
					Provider:             "azure",
					FirstTimestamp:       time.Now().UTC().Add(-2 * 24 * time.Hour),
					NamespaceDescription: "mynamespace3",
				},
				// 5 days old cluster
				{
					Installation:         "ginger",
					ID:                   "abc12",
					Release:              "1.2.3-myversion",
					Provider:             "aws",
					FirstTimestamp:       time.Now().UTC().Add(-5 * 24 * time.Hour),
					NamespaceDescription: "<unknown>",
				},
			},
			errors:     []error{errors.New("First error"), errors.New("Second error")},
			goldenFile: "success.golden",
		},
		{
			name:       "Errors-only",
			clusters:   []cortex.Cluster{},
			errors:     []error{errors.New("nothing but failure")},
			goldenFile: "error.golden",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			expectedBytes, err := fs.ReadFile("golden/" + tt.goldenFile)
			if err != nil {
				t.Fatalf("could not open golden file - %s", err)
			}

			expected := string(expectedBytes)
			report, err := Render(tt.clusters, tt.errors)
			if err != nil {
				t.Fatalf("expected err to be nil, got %s", err)
			}

			if !cmp.Equal(expected, report) {
				t.Fatalf("report doesn't match expected: %s", cmp.Diff(expected, report))
			}
		})
	}
}
