{{/* vim: set filetype=mustache: */}}
{{- define "resourcePolice.containerImage" -}}
{{- if .Values.image.tag }}
{{- .Values.image.registry }}/{{ .Values.image.name }}:{{ .Values.image.tag }}
{{- else }}
{{- .Values.image.registry }}/{{ .Values.image.name }}:{{ .Chart.Version }}
{{- end -}}
{{- end -}}
