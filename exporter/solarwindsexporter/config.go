package solarwindsexporter

import (
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

type Config struct {
	configmodels.ExporterSettings  `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.TimeoutSettings `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
	exporterhelper.QueueSettings   `mapstructure:"sending_queue"`
	exporterhelper.RetrySettings   `mapstructure:"retry_on_failure"`
	TraceEndpoint                  string                   `mapstructure:"trace_endpoint"`
	MetricEndpoint                 string                   `mapstructure:"metric_endpoint"`
	LogEndpoint                    string                   `mapstructure:"log_endpoint"`
	ServiceKey                     string                   `mapstructure:"service_key"`
	configgrpc.GRPCClientSettings  `mapstructure:",squash"` // squash ensures fields are correctly decoded in embedded struct.
}
