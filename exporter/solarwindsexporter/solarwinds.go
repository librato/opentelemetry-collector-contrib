package solarwindsexporter

import (
	"context"
	"errors"

	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/consumer/pdata"
)

type exporterImp struct {
	// Input configuration.
	config *Config
	//	w      *grpcSender
}

var (
	errPermanentError = consumererror.Permanent(errors.New("fatal error sending to server"))
)

// Crete new exporter and start it. The exporter will begin connecting but
// this function may return before the connection is established.
func newExporter(cfg configmodels.Exporter) (*exporterImp, error) {
	oCfg := cfg.(*Config)

	if oCfg.Endpoint == "" {
		return nil, errors.New("OTLP exporter config requires an Endpoint")
	}

	e := &exporterImp{}
	e.config = oCfg
	// w, err := newGrpcSender(oCfg)
	// if err != nil {
	// 	return nil, err
	// }
	// e.w = w
	return e, nil
}

func (e *exporterImp) shutdown(context.Context) error {
	//	return e.w.stop()
	return nil
}

func (e *exporterImp) pushTraceData(ctx context.Context, td pdata.Traces) (int, error) {
	// request := &otlptrace.ExportTraceServiceRequest{
	// 	ResourceSpans: pdata.TracesToOtlp(td),
	// }
	// err := e.w.exportTrace(ctx, request)

	// if err != nil {
	// 	return td.SpanCount(), fmt.Errorf("failed to push trace data via OTLP exporter: %w", err)
	// }
	return 0, nil
}

func (e *exporterImp) pushMetricsData(ctx context.Context, md pdata.Metrics) (int, error) {
	// request := &otlpmetrics.ExportMetricsServiceRequest{
	// 	ResourceMetrics: pdata.MetricsToOtlp(md),
	// }
	// err := e.w.exportMetrics(ctx, request)

	// if err != nil {
	// 	return md.MetricCount(), fmt.Errorf("failed to push metrics data via OTLP exporter: %w", err)
	// }
	return 0, nil
}

func (e *exporterImp) pushLogData(ctx context.Context, logs pdata.Logs) (int, error) {
	// request := &otlplogs.ExportLogsServiceRequest{
	// 	ResourceLogs: internal.LogsToOtlp(logs.InternalRep()),
	// }
	// err := e.w.exportLogs(ctx, request)

	// if err != nil {
	// 	return logs.LogRecordCount(), fmt.Errorf("failed to push log data via OTLP exporter: %w", err)
	// }
	return 0, nil
}
