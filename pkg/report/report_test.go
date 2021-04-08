package report

import (
	"embed"
	"testing"

	"github.com/google/go-cmp/cmp"
)

//go:embed report.golden
var fs embed.FS

func Test_RenderReport(t *testing.T) {
	expectedBytes, _ := fs.ReadFile("report.golden")
	expected := string(expectedBytes)

	report, err := Render(
		[]string{"gaia/def34", "ginger/abc12"},
		[]string{"First error", "Second error"})
	if err != nil {
		t.Fatalf("expected err to be nil, got %s", err)
	}
	if !cmp.Equal(expected, report) {
		t.Fatalf("report doesn't match expected: %s", cmp.Diff(expected, report))
	}
}
