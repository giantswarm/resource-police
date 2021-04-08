// Package report handles the rendering of the
// HTML report that is then sent out via Slack.
package report

import (
	"bytes"
	"embed"
	"fmt"
	"strings"
	"text/template"

	"github.com/giantswarm/microerror"
)

//go:embed template.gotmpl
var f embed.FS

type Cluster struct {
	ID               string
	InstallationName string
}

type TemplateData struct {
	Clusters []Cluster
	Errors   []string
}

func Render(clusters []string, errors []string) (string, error) {
	fmt.Println("Rendering report")

	reportTemplate, _ := f.ReadFile("template.gotmpl")

	t := template.Must(template.New("report-template").Parse(string(reportTemplate)))

	myData := TemplateData{
		Clusters: []Cluster{},
		Errors:   errors,
	}

	for _, c := range clusters {
		parts := strings.Split(c, "/")
		myData.Clusters = append(myData.Clusters, Cluster{
			ID:               parts[1],
			InstallationName: parts[0],
		})
	}

	var renderedReport bytes.Buffer
	err := t.Execute(&renderedReport, myData)
	if err != nil {
		return "", microerror.Mask(err)
	}

	fmt.Println("Report has been rendered")

	return renderedReport.String(), nil
}
