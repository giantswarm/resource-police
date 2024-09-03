// Package report handles the rendering of the
// HTML report that is then sent out via Slack.
package report

import (
	"bytes"
	"embed"
	"log"
	"sort"
	"strings"
	"text/template"
	"time"

	"github.com/giantswarm/microerror"

	"github.com/giantswarm/resource-police/pkg/cortex"
)

//go:embed template.gotmpl
var f embed.FS

type TemplateData struct {
	Clusters3Days  []cortex.Cluster
	Clusters1Day   []cortex.Cluster
	Clusters3Hours []cortex.Cluster
	ClustersRest   []cortex.Cluster
	Errors         []error
}

func Render(clusters []cortex.Cluster, errors []error) (string, error) {
	log.Println("Rendering report")

	reportTemplate, _ := f.ReadFile("template.gotmpl")

	t := template.Must(template.New("report-template").Parse(string(reportTemplate)))

	myData := TemplateData{
		Errors: errors,
	}

	// Sort by provider and installation name
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Installation < clusters[j].Installation
	})

	// Afterwards sort by provider
	sort.Slice(clusters, func(i, j int) bool {
		return clusters[i].Provider < clusters[j].Provider
	})

	// Group clusters by age
	// (more than 3 days, more than 1 day, more than 3 hours, rest).
	now := time.Now().UTC()
	for i := 0; i < len(clusters); i++ {
		switch age := now.Sub(clusters[i].FirstTimestamp); {
		case age > 3*24*time.Hour:
			myData.Clusters3Days = append(myData.Clusters3Days, clusters[i])
		case age > 24*time.Hour:
			myData.Clusters1Day = append(myData.Clusters1Day, clusters[i])
		case age > 3*time.Hour:
			myData.Clusters3Hours = append(myData.Clusters3Hours, clusters[i])
		default:
			myData.ClustersRest = append(myData.ClustersRest, clusters[i])
		}
	}

	var renderedReport bytes.Buffer
	err := t.Execute(&renderedReport, myData)
	if err != nil {
		return "", microerror.Mask(err)
	}

	log.Println("Report has been rendered")

	return strings.TrimSpace(renderedReport.String()) + "\n", nil
}
