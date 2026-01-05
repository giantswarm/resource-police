{{/* vim: set filetype=mustache: */}}
{{- define "resourcePolice.containerImage" -}}
{{- if .Values.image.tag }}
{{- .Values.image.registry }}/{{ .Values.image.name }}:{{ .Values.image.tag }}
{{- else }}
{{- .Values.image.registry }}/{{ .Values.image.name }}:{{ .Chart.AppVersion }}
{{- end -}}
{{- end -}}

{{/*
Common labels
*/}}
{{- define "labels.common" -}}
app: {{ .Values.name }}
giantswarm.io/service-type: "managed"
application.giantswarm.io/team: {{ index .Chart.Annotations "io.giantswarm.application.team" | quote }}
{{- end -}}
