{{/*
Expand the name of the chart.
*/}}
{{- define "prometheus-isilon-exporter.name" -}}
{{- default .Chart.Name .Values.nameOverride | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Create a default fully qualified app name.
We truncate at 63 chars because some Kubernetes name fields are limited to this (by the DNS naming spec).
If release name contains chart name it will be used as a full name.
*/}}
{{- define "prometheus-isilon-exporter.fullname" -}}
{{- if .Values.fullnameOverride }}
{{- .Values.fullnameOverride | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- $name := default .Chart.Name .Values.nameOverride }}
{{- if contains $name .Release.Name }}
{{- .Release.Name | trunc 63 | trimSuffix "-" }}
{{- else }}
{{- printf "%s-%s" .Release.Name $name | trunc 63 | trimSuffix "-" }}
{{- end }}
{{- end }}
{{- end }}

{{/*
Create chart name and version as used by the chart label.
*/}}
{{- define "prometheus-isilon-exporter.chart" -}}
{{- printf "%s-%s" .Chart.Name .Chart.Version | replace "+" "_" | trunc 63 | trimSuffix "-" }}
{{- end }}

{{/*
Common labels
*/}}
{{- define "prometheus-isilon-exporter.labels" -}}
helm.sh/chart: {{ include "prometheus-isilon-exporter.chart" . }}
{{ include "prometheus-isilon-exporter.selectorLabels" . }}
{{- if .Chart.AppVersion }}
app.kubernetes.io/version: {{ .Chart.AppVersion | quote }}
{{- end }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
{{- end }}

{{/*
Selector labels
*/}}
{{- define "prometheus-isilon-exporter.selectorLabels" -}}
app.kubernetes.io/name: {{ include "prometheus-isilon-exporter.name" . }}
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end }}

{{/*
Create the name of the service account to use
*/}}
{{- define "prometheus-isilon-exporter.serviceAccountName" -}}
{{- if .Values.serviceAccount.create }}
{{- default (include "prometheus-isilon-exporter.fullname" .) .Values.serviceAccount.name }}
{{- else }}
{{- default "default" .Values.serviceAccount.name }}
{{- end }}
{{- end }}

# {{/*
# Add normal and secret env variables 
# */}}
# {{- define "helpers.list-env-variables"}}
# {{- range $key, $val := .Values.env.secret }}
# - name: {{ $key }}
#   valueFrom:
#     secretKeyRef:
#       name: app-env-vars-secret
#       key: {{ $key }}
# {{- end}}
# {{- range $key, $val := .Values.env.normal }}
# - name: {{ $key }}
#   valueFrom:
#     configMapKeyRef:
#       name: app-env-vars-configmap
#       key: {{ $key }}
# {{- end}}
# {{- end }}
