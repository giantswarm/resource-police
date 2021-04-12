package report

import (
	"embed"
	"errors"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed golden/*
var fs embed.FS

func Test_RenderReport(t *testing.T) {
	tests := []struct {
		name       string
		clusters   []string
		errors     []error
		goldenFile string
	}{
		{
			name:       "Success",
			clusters:   []string{"gaia/def34", "ginger/abc12"},
			errors:     []error{errors.New("First error"), errors.New("Second error")},
			goldenFile: "success.golden",
		},
		{
			name:       "Errors-only",
			clusters:   []string{},
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
