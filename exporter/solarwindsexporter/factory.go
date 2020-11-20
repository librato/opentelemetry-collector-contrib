package solarwindsexporter

import (
	"context"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/exporter/exporterhelper"
)

const (
	typeStr         = "solarwinds"
	defaultEndpoint = "https://dc.services.visualstudio.com/v2/track"
)

// NewFactory creates a factory for OTLP exporter.
func NewFactory() component.ExporterFactory {
	return exporterhelper.NewFactory(
		typeStr,
		createDefaultConfig,
		exporterhelper.WithTraces(createTraceExporter),
		exporterhelper.WithMetrics(createMetricsExporter),
		exporterhelper.WithLogs(createLogsExporter))
}

func createDefaultConfig() configmodels.Exporter {
	return &Config{
		ExporterSettings: configmodels.ExporterSettings{
			TypeVal: typeStr,
			NameVal: typeStr,
		},
		TimeoutSettings: exporterhelper.CreateDefaultTimeoutSettings(),
		RetrySettings:   exporterhelper.CreateDefaultRetrySettings(),
		QueueSettings:   exporterhelper.CreateDefaultQueueSettings(),
		GRPCClientSettings: configgrpc.GRPCClientSettings{
			Headers: map[string]string{},
			// We almost read 0 bytes, so no need to tune ReadBufferSize.
			WriteBufferSize: 512 * 1024,
		},
	}
}

func createTraceExporter(
	_ context.Context,
	_ component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.TraceExporter, error) {
	oce, err := newExporter(cfg)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)
	oexp, err := exporterhelper.NewTraceExporter(
		cfg,
		oce.pushTraceData,
		exporterhelper.WithTimeout(oCfg.TimeoutSettings),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithShutdown(oce.shutdown))
	if err != nil {
		return nil, err
	}

	return oexp, nil
}

func createMetricsExporter(
	_ context.Context,
	_ component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.MetricsExporter, error) {
	oce, err := newExporter(cfg)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)
	oexp, err := exporterhelper.NewMetricsExporter(
		cfg,
		oce.pushMetricsData,
		exporterhelper.WithTimeout(oCfg.TimeoutSettings),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithShutdown(oce.shutdown),
	)
	if err != nil {
		return nil, err
	}

	return oexp, nil
}

func createLogsExporter(
	_ context.Context,
	_ component.ExporterCreateParams,
	cfg configmodels.Exporter,
) (component.LogsExporter, error) {
	oce, err := newExporter(cfg)
	if err != nil {
		return nil, err
	}
	oCfg := cfg.(*Config)
	oexp, err := exporterhelper.NewLogsExporter(
		cfg,
		oce.pushLogData,
		exporterhelper.WithTimeout(oCfg.TimeoutSettings),
		exporterhelper.WithRetry(oCfg.RetrySettings),
		exporterhelper.WithQueue(oCfg.QueueSettings),
		exporterhelper.WithShutdown(oce.shutdown),
	)
	if err != nil {
		return nil, err
	}

	return oexp, nil
}
