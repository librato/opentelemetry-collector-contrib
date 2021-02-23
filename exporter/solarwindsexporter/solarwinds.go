package solarwindsexporter

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/appoptics/appoptics-apm-go/v1/ao"
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

	// if oCfg.Endpoint == "" {
	// 	return nil, errors.New("OTLP exporter config requires an Endpoint")
	// }

	e := &exporterImp{}
	e.config = oCfg
	// w, err := newGrpcSender(oCfg)
	// if err != nil {
	// 	return nil, err
	// }
	// e.w = w
	return e, nil
}

const (
	xtraceVersionHeader = "2B"
	sampledFlags        = "01"
)

func toXTraceID(otTraceId pdata.TraceID, otSpanId pdata.SpanID) string {
	taskId := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%0-40v", otTraceId.HexString()), " ", "0"))
	opId := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%0-16v", otSpanId.HexString()), " ", "0"))
	return xtraceVersionHeader + taskId + opId + sampledFlags
}

func (e *exporterImp) shutdown(context.Context) error {
	//	return e.w.stop()
	return nil
}

func (e *exporterImp) pushTraceData(ctx context.Context, td pdata.Traces) (int, error) {

	spans := td.ResourceSpans()
	for i := 0; i < spans.Len(); i++ {
		libSpans := spans.At(i).InstrumentationLibrarySpans()
		processedIDs := make(map[string]struct{})
		for j := 0; j < libSpans.Len(); j++ {
			libSpan := libSpans.At(i).Spans()
			var traceContext context.Context
			for k := 0; k < libSpan.Len(); k++ {
				span := libSpan.At(k)
				xTraceID := toXTraceID(span.TraceID(), span.SpanID())

				if _, ok := processedIDs[xTraceID]; !ok {
					fmt.Printf("XTrace ID %v existing processed %v i/j/k %v %v %v \n ", xTraceID, processedIDs, i, j, k)
					processedIDs[xTraceID] = struct{}{}
					startOverrides := ao.Overrides{
						ExplicitTS:    time.Unix(0, (int64)(span.StartTime())),
						ExplicitMdStr: xTraceID,
					}
					endOverrides := ao.Overrides{
						ExplicitTS: time.Unix(0, (int64)(span.EndTime())),
					}
					if len(span.ParentSpanID().Bytes()) == 0 {
						fmt.Println("Trace start===============")
						trace := ao.NewTraceWithOverrides(span.Name(), startOverrides, nil)
						traceContext = ao.NewContext(context.Background(), trace)
						trace.SetStartTime(time.Unix(0, (int64)(span.StartTime()))) //this is for histogram only
						fmt.Println("Trace end===============")
						trace.EndWithOverrides(endOverrides)
						fmt.Printf("Root span %v with start overrides : %+v and end overrides : %+v\n", span.Name(), startOverrides, endOverrides)
					} else {
						//parentXTraceID := toXTraceID(span.TraceID(), span.ParentSpanID())
						aoSpan, _ := ao.BeginSpanWithOverrides(traceContext, span.Name(), ao.SpanOptions{}, startOverrides)
						aoSpan.EndWithOverrides(endOverrides)
						fmt.Printf("child span %v with context %+v and start overrides : %+v and end overrides : %+v\n", span.Name(), traceContext, startOverrides, endOverrides)
					}
				}
			}
		}
	}

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
