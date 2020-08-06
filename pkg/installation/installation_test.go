package installation

import (
	"testing"

	"github.com/google/go-cmp/cmp"
)

const expectedRenderedReport = `*Test clusters that should be deleted*

- ` + "`ginger` / `abc12`" + ` - cluster (1d old) - ping @creator
- ` + "`gaia` / `def34`" + ` - other (2d old)

Please check <https://intranet.giantswarm.io/docs/dev-and-releng/test-environments/|our policy> on how to keep test clusters alive.
`

func Test_RenderReport(t *testing.T) {
	report, err := RenderReport([]*Cluster{
		{
			ID:               "abc12",
			Name:             "cluster",
			AgeString:        "1d",
			Creator:          "creator",
			InstallationName: "ginger",
		},
		{
			ID:               "def34",
			Name:             "other",
			AgeString:        "2d",
			InstallationName: "gaia",
		},
	}, []string{"An error occurred"})
	if err != nil {
		t.Fatalf("expected err to be nil, got %s", err)
	}
	if !cmp.Equal(expectedRenderedReport, report) {
		t.Fatalf("report doesn't match expected: %s", cmp.Diff(expectedRenderedReport, report))
	}
}
