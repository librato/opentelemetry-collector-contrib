module github.com/librato/opentelemetry-collector-contrib/exporter/solarwindsexporter

go 1.15

require (
	github.com/appoptics/appoptics-apm-go v1.14.0
	github.com/stretchr/testify v1.7.0
	go.opentelemetry.io/collector v0.13.0
	go.uber.org/zap v1.16.0
)

replace github.com/appoptics/appoptics-apm-go => ../../../appoptics-apm-go
