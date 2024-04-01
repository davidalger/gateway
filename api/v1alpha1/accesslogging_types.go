// Copyright Envoy Gateway Authors
// SPDX-License-Identifier: Apache-2.0
// The full text of the Apache license is available in the LICENSE file at
// the root of the repo.

package v1alpha1

import gwapiv1 "sigs.k8s.io/gateway-api/apis/v1"

type ProxyAccessLog struct {
	// Disable disables access logging for managed proxies if set to true.
	Disable bool `json:"disable,omitempty"`
	// Settings defines accesslog settings for managed proxies.
	// If unspecified, will send default format to stdout.
	// +optional
	Settings []ProxyAccessLogSetting `json:"settings,omitempty"`
}

type ProxyAccessLogSetting struct {
	// Format defines the format of accesslog.
	Format ProxyAccessLogFormat `json:"format"`
	// Sinks defines the sinks of accesslog.
	// +kubebuilder:validation:MinItems=1
	Sinks []ProxyAccessLogSink `json:"sinks"`
}

type ProxyAccessLogFormatType string

const (
	// ProxyAccessLogFormatTypeText defines the text accesslog format.
	ProxyAccessLogFormatTypeText ProxyAccessLogFormatType = "Text"
	// ProxyAccessLogFormatTypeJSON defines the JSON accesslog format.
	ProxyAccessLogFormatTypeJSON ProxyAccessLogFormatType = "JSON"
	// TODO: support format type "mix" in the future.
)

// ProxyAccessLogFormat defines the format of accesslog.
// By default accesslogs are written to standard output.
// +union
//
// +kubebuilder:validation:XValidation:rule="self.type == 'Text' ? has(self.text) : !has(self.text)",message="If AccessLogFormat type is Text, text field needs to be set."
// +kubebuilder:validation:XValidation:rule="self.type == 'JSON' ? has(self.json) : !has(self.json)",message="If AccessLogFormat type is JSON, json field needs to be set."
type ProxyAccessLogFormat struct {
	// Type defines the type of accesslog format.
	// +kubebuilder:validation:Enum=Text;JSON
	// +unionDiscriminator
	Type ProxyAccessLogFormatType `json:"type,omitempty"`
	// Text defines the text accesslog format, following Envoy accesslog formatting,
	// It's required when the format type is "Text".
	// Envoy [command operators](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#command-operators) may be used in the format.
	// The [format string documentation](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#config-access-log-format-strings) provides more information.
	// +optional
	Text *string `json:"text,omitempty"`
	// JSON is additional attributes that describe the specific event occurrence.
	// Structured format for the envoy access logs. Envoy [command operators](https://www.envoyproxy.io/docs/envoy/latest/configuration/observability/access_log/usage#command-operators)
	// can be used as values for fields within the Struct.
	// It's required when the format type is "JSON".
	// +optional
	JSON map[string]string `json:"json,omitempty"`
}

type ProxyAccessLogSinkType string

const (
	// ProxyAccessLogSinkTypeALS defines the gRPC Access Log Service (ALS) sink.
	// The service must implement the Envoy gRPC Access Log Service streaming API:
	// https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/accesslog/v3/als.proto
	ProxyAccessLogSinkTypeALS ProxyAccessLogSinkType = "ALS"
	// ProxyAccessLogSinkTypeFile defines the file accesslog sink.
	ProxyAccessLogSinkTypeFile ProxyAccessLogSinkType = "File"
	// ProxyAccessLogSinkTypeOpenTelemetry defines the OpenTelemetry accesslog sink.
	// When the provider is Kubernetes, EnvoyGateway always sends `k8s.namespace.name`
	// and `k8s.pod.name` as additional attributes.
	ProxyAccessLogSinkTypeOpenTelemetry ProxyAccessLogSinkType = "OpenTelemetry"
)

// ProxyAccessLogSink defines the sink of accesslog.
// +union
//
// +kubebuilder:validation:XValidation:rule="self.type == 'ALS' ? has(self.als) : !has(self.als)",message="If AccessLogSink type is ALS, als field needs to be set."
// +kubebuilder:validation:XValidation:rule="self.type == 'File' ? has(self.file) : !has(self.file)",message="If AccessLogSink type is File, file field needs to be set."
// +kubebuilder:validation:XValidation:rule="self.type == 'OpenTelemetry' ? has(self.openTelemetry) : !has(self.openTelemetry)",message="If AccessLogSink type is OpenTelemetry, openTelemetry field needs to be set."
type ProxyAccessLogSink struct {
	// Type defines the type of accesslog sink.
	// +kubebuilder:validation:Enum=ALS;File;OpenTelemetry
	// +unionDiscriminator
	Type ProxyAccessLogSinkType `json:"type,omitempty"`
	// ALS defines the gRPC Access Log Service (ALS) sink.
	// +optional
	ALS *ALSEnvoyProxyAccessLog `json:"als,omitempty"`
	// File defines the file accesslog sink.
	// +optional
	File *FileEnvoyProxyAccessLog `json:"file,omitempty"`
	// OpenTelemetry defines the OpenTelemetry accesslog sink.
	// +optional
	OpenTelemetry *OpenTelemetryEnvoyProxyAccessLog `json:"openTelemetry,omitempty"`
}

type ALSEnvoyProxyAccessLogType string

const (
	// ALSEnvoyProxyAccessLogTypeHTTP defines the HTTP access log type and will populate StreamAccessLogsMessage.http_logs.
	ALSEnvoyProxyAccessLogTypeHTTP ALSEnvoyProxyAccessLogType = "HTTP"
	// ALSEnvoyProxyAccessLogTypeTCP defines the TCP access log type and will populate StreamAccessLogsMessage.tcp_logs.
	ALSEnvoyProxyAccessLogTypeTCP ALSEnvoyProxyAccessLogType = "TCP"
)

// ALSEnvoyProxyAccessLog defines the gRPC Access Log Service (ALS) sink.
// The service must implement the Envoy gRPC Access Log Service streaming API:
// https://www.envoyproxy.io/docs/envoy/latest/api-v3/service/accesslog/v3/als.proto
// Access log format information is passed in the form of gRPC metadata when the
// stream is established. Specifically, the following metadata is passed:
//
// - `x-accesslog-text` - The access log format string when a Text format is used.
//
// - `x-accesslog-attr` - JSON encoded key/value pairs when a JSON format is used.
//
// +kubebuilder:validation:XValidation:message="BackendRef only supports Service Kind.",rule="!has(self.backendRef.kind) || self.backendRef.kind == 'Service'"
// +kubebuilder:validation:XValidation:rule="self.type == 'HTTP' || !has(self.http)",message="The http field may only be set when type is HTTP."
type ALSEnvoyProxyAccessLog struct {
	// BackendRef references a Kubernetes object that represents the gRPC service to which
	// the access logs will be sent. Currently only Service is supported.
	BackendRef gwapiv1.BackendObjectReference `json:"backendRef,omitempty"`
	// LogName defines the friendly name of the access log to be returned in
	// StreamAccessLogsMessage.Identifier. This allows the access log server
	// to differentiate between different access logs coming from the same Envoy.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:default="accesslog"
	LogName string `json:"logName,omitempty"`
	// Type defines the type of accesslog. Supported types are "HTTP" and "TCP". Defaults to "HTTP" when not specified.
	// +kubebuilder:validation:Enum=HTTP;TCP
	// +unionDiscriminator
	// +kubebuilder:default="HTTP"
	Type ALSEnvoyProxyAccessLogType `json:"type,omitempty"`
	// HTTP defines additional configuration specific to HTTP access logs.
	// +optional
	HTTP *ALSEnvoyProxyHTTPAccessLogConfig `json:"http,omitempty"`
}

type ALSEnvoyProxyHTTPAccessLogConfig struct {
	// RequestHeaders defines request headers to include in log entries sent to the access log service.
	// +optional
	RequestHeaders []string `json:"requestHeaders,omitempty" yaml:"requestHeaders,omitempty"`
	// ResponseHeaders defines response headers to include in log entries sent to the access log service.
	// +optional
	ResponseHeaders []string `json:"responseHeaders,omitempty" yaml:"responseHeaders,omitempty"`
	// ResponseTrailers defines response trailers to include in log entries sent to the access log service.
	// +optional
	ResponseTrailers []string `json:"responseTrailers,omitempty" yaml:"responseTrailers,omitempty"`
}

type FileEnvoyProxyAccessLog struct {
	// Path defines the file path used to expose envoy access log(e.g. /dev/stdout).
	// +kubebuilder:validation:MinLength=1
	Path string `json:"path,omitempty"`
}

// TODO: consider reuse ExtensionService?
type OpenTelemetryEnvoyProxyAccessLog struct {
	// Host define the extension service hostname.
	Host string `json:"host"`
	// Port defines the port the extension service is exposed on.
	//
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:default=4317
	Port int32 `json:"port,omitempty"`
	// Resources is a set of labels that describe the source of a log entry, including envoy node info.
	// It's recommended to follow [semantic conventions](https://opentelemetry.io/docs/reference/specification/resource/semantic_conventions/).
	// +optional
	Resources map[string]string `json:"resources,omitempty"`

	// TODO: support more OpenTelemetry accesslog options(e.g. TLS, auth etc.) in the future.
}
