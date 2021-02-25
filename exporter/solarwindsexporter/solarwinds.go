package solarwindsexporter

import (
	"context"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/appoptics/appoptics-apm-go/v1/ao"
	v1 "github.com/census-instrumentation/opencensus-proto/gen-go/trace/v1"
	"go.opentelemetry.io/collector/config/configmodels"
	"go.opentelemetry.io/collector/consumer/consumererror"
	"go.opentelemetry.io/collector/consumer/pdata"
	"go.opentelemetry.io/collector/translator/internaldata"
	"strings"
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

	//os.Setenv("APPOPTICS_SERVICE_KEY", oCfg.ServiceKey) //this is too late, the Config is already loaded for agent
	//fmt.Printf("Setting APPOPTICS_SERVICE_KEY to %v all values %+v", oCfg.ServiceKey, oCfg)
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

var wsKeyMap = map[string]string{
	"http.method":      "HTTPMethod",
	"http.url":         "URL",
	"http.status_code": "Status",
}
var queryKeyMap = map[string]string{
	"db.connection_string": "RemoteHost",
	"db.name":              "Database",
	"db.statement":         "Query",
	"db.system":            "Flavor",
}

//func toXTraceID(otTraceId pdata.TraceID, otSpanId pdata.SpanID) string {
//	taskId := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%0-40v", otTraceId.HexString()), " ", "0"))
//	opId := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%0-16v", otSpanId.HexString()), " ", "0"))
//	return xtraceVersionHeader + taskId + opId + sampledFlags
//}

func getXTraceID(traceID []byte, spanID []byte) string {
	taskId := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%0-40v", hex.EncodeToString(traceID)), " ", "0"))
	opId := strings.ToUpper(strings.ReplaceAll(fmt.Sprintf("%0-16v", hex.EncodeToString(spanID)), " ", "0"))
	return xtraceVersionHeader + taskId + opId + sampledFlags
}

func extractWebserverKvs(span *v1.Span) []interface{} {
	return extractSpecKvs(span, wsKeyMap, "ws")
}

func extractQueryKvs(span *v1.Span) []interface{} {
	return extractSpecKvs(span, queryKeyMap, "query")
}

func extractSpecKvs(span *v1.Span, lookup map[string]string, specValue string) []interface{} {
	attrMap := span.Attributes.AttributeMap
	result := []interface{}{}
	for otKey, aoKey := range lookup {
		if val, ok := attrMap[otKey]; ok {
			result = append(result, aoKey)
			result = append(result, fromAttributeValue(val))
		}
	}
	if len(result) > 0 {
		result = append(result, "Spec")
		result = append(result, specValue)
	}
	return result
}

func fromAttributeValue(attributeValue *v1.AttributeValue) interface{} {
	switch attributeValue.GetValue().(type) {
	case *v1.AttributeValue_StringValue:
		return attributeValue.GetStringValue().Value
	case *v1.AttributeValue_IntValue:
		return attributeValue.GetIntValue()
	case *v1.AttributeValue_DoubleValue:
		return attributeValue.GetDoubleValue()
	case *v1.AttributeValue_BoolValue:
		return attributeValue.GetBoolValue()
	default:
		return nil
	}
}

func extractKvs(span *v1.Span) []interface{} {
	var kvs []interface{}
	for key, attributeValue := range span.Attributes.AttributeMap {
		kvs = append(kvs, key)
		kvs = append(kvs, fromAttributeValue(attributeValue))
	}
	if len(span.ParentSpanId) == 0 { //root span, attempt to extract webserver KVs
		kvs = append(kvs, extractWebserverKvs(span)...)
	}
	kvs = append(kvs, extractQueryKvs(span)...)

	return kvs
}

func (e *exporterImp) shutdown(context.Context) error {
	//	return e.w.stop()
	return nil
}

func (e *exporterImp) pushTraceData(ctx context.Context, td pdata.Traces) (int, error) {
	octds := internaldata.TraceDataToOC(td)
	for _, octd := range octds {
		for _, span := range octd.Spans {
			xTraceID := getXTraceID(span.TraceId, span.SpanId)

			startOverrides := ao.Overrides{
				ExplicitTS:    span.StartTime.AsTime(),
				ExplicitMdStr: xTraceID,
			}
			endOverrides := ao.Overrides{
				ExplicitTS: span.EndTime.AsTime(),
			}
			var traceContext context.Context
			kvs := extractKvs(span)

			if len(span.ParentSpanId) == 0 {
				trace := ao.NewTraceWithOverrides(span.Name.Value, startOverrides, nil)
				traceContext = ao.NewContext(context.Background(), trace)
				trace.SetStartTime(span.StartTime.AsTime()) //this is for histogram only
				trace.EndWithOverrides(endOverrides, kvs...)
			} else {
				parentXTraceID := getXTraceID(span.TraceId, span.ParentSpanId)
				traceContext = ao.FromXTraceIDContext(context.Background(), parentXTraceID)
				aoSpan, _ := ao.BeginSpanWithOverrides(traceContext, span.Name.Value, ao.SpanOptions{}, startOverrides)
				aoSpan.EndWithOverrides(endOverrides, kvs...)
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
