{{- define "synthkit.fullname" -}}
{{- if eq .Release.Name .Chart.Name -}}
{{- .Release.Name | trunc 63 | trimSuffix "-" -}}
{{- else -}}
{{- printf "%s-%s" .Release.Name .Chart.Name | trunc 63 | trimSuffix "-" -}}
{{- end -}}
{{- end -}}

{{- define "synthkit.selectorLabels" -}}
app.kubernetes.io/name: synthkit
app.kubernetes.io/instance: {{ .Release.Name }}
{{- end -}}

{{- define "synthkit.labels" -}}
{{ include "synthkit.selectorLabels" . }}
app.kubernetes.io/managed-by: {{ .Release.Service }}
helm.sh/chart: {{ printf "%s-%s" .Chart.Name .Chart.Version | quote }}
{{- end -}}

{{- define "synthkit.claimName" -}}
{{- default (printf "%s-state" (include "synthkit.fullname" . | trunc 57 | trimSuffix "-")) .Values.persistence.existingClaim -}}
{{- end -}}
