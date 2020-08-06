package installation

// Backticks in multi-line strings can't be escaped, so we're appending a normal double quoted string containing
// backticks in the middle of this template.
const reportTemplate = `*Test clusters that should be deleted*

{{ range  .Clusters -}}
- ` + "`{{ .InstallationName }}` / `{{ .ID }}`" + ` - {{ .Name }} ({{ .AgeString }} old){{ with .Creator }} - ping @{{ . }}{{ end }}
{{ end }}
Please check <https://intranet.giantswarm.io/docs/dev-and-releng/test-environments/|our policy> on how to keep test clusters alive.

{{ if .Errors -}}
By the way, some errors occurred:

{{ range  .Errors -}}
- {{ . }}
{{ end }}
Please check the resource-police configuration and installation IP whitelists to ensure
that all installations are accessible by resource-police.
{{ end }}
`
