To start the collector with our span exporter. Simply use `docker-compose up` in this current directory.

The docker container will bind to 0.0.0.0:4317 of the docker host, make sure the port is made available. On EC2 instance port might be taken by the AWS otel process.
`/opt/aws/aws-otel-collector/bin/aws-otel-collector --config /opt/aws/aws-otel-collector/etc/config.yaml` . You might either need to kill it or bind the docker instance to a different port

The config file `otel-config.yaml` is mounted to the running docker-compose instance

### Remarks
The DockerFile first run `RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder`, take note that we have to use the one in github.com/observatorium/opentelemetry-collector-builder instead of the newer github.com/open-telemetry/opentelemetry-collector-builder as the newer one contains API that does not match our implementation (which was based on v0.13.1 of opentelemetry-collector-contrib)

If the DockerFile uses `FROM golang:1.16-buster` (golang 1.16), then the builder from `observatorium` seems to have issue to resolve the dependencies 
```
Error: failed to compile the OpenTelemetry Collector distribution: exit status 1. Output: "go: go.opentelemetry.io/collector@v0.13.1-0.20201020175630-99cb5b244aad: missing go.sum entry; to add it:\n\tgo mod download go.opentelemetry.io/collector\n"
```
Using the builder from `open-telemetry` resolves the dependencies fine, but have other errors (probably due to API mismatch) such as:
```
Error: failed to compile the OpenTelemetry Collector distribution: exit status 2. Output: "# github.com/jpkroehling/opentelemetry-collector-builder\n./components.go:83:20: undefined: consumererror.Combine\n./main.go:34:10: undefined: component.BuildInfo\n./main.go:40:26: undefined: service.CollectorSettings\n"
```

So the solution of using `FROM golang:1.16-buster` is rather ugly - First `RUN opentelemetry-collector-builder ; return 0` on the `open-telemetry` builder (`RUN GO111MODULE=on go get github.com/open-telemetry/opentelemetry-collector-builder`), it will fail compilation but at least pull the correct dependencies, then again we `RUN opentelemetry-collector-builder` but on the `observatorium` builder (by doing `RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder`)

It works, but very ugly

Later on, it's found that we can use `observatorium` builder with `FROM golang:1.15-buster`, which pulls dependencies correctly!


