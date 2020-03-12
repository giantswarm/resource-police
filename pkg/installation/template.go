package installation

const ReportTemplate = `*Report about clusters older than 8 hours*
{{ range $i, $installation := . }}
Installation: *{{ $installation.Name }}*
{{ if $installation.Clusters -}}
Clusters:
{{ range $j, $cluster := $installation.Clusters -}}
  - id: *{{ $cluster.ID }}* | name: *{{ $cluster.Name}}* | age: *{{ $cluster.Age }}*
{{ end -}}
{{ end -}}
{{ end -}}
`
