FROM golang:1.16.5-buster
 
#RUN cd ~ && git clone https://github.com/open-telemetry/opentelemetry-collector-builder.git && pwd && cd opentelemetry-collector-builder && pwd && ls && git checkout tags/v0.29.0 && echo "here" && ls

RUN GO111MODULE=on go get github.com/open-telemetry/opentelemetry-collector-builder
#RUN GO111MODULE=on go get github.com/open-telemetry/opentelemetry-collector-builder@v0.25.0
#RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder
#RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder@v0.1.2

RUN git clone https://github.com/librato/opentelemetry-collector-contrib.git && cd opentelemetry-collector-contrib && git checkout AO-18760-ao-exporter && cd ..

RUN git clone https://github.com/appoptics/appoptics-apm-go.git && cd appoptics-apm-go && git checkout AO-18760-ao-exporter && cd ..

COPY .otelcol-builder.yaml /root/.otelcol-builder.yaml

## This build would fail, but at least it would pull go.opentelemetry.io/collector@0.13.1-0.20201020175630-99cb5b244aad correctly
RUN opentelemetry-collector-builder; exit 0

RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder

## This requires the previous step to pull go.opentelemetry.io/collector@0.13.1-0.20201020175630-99cb5b244aad, otherwise it cannot resolve it
RUN opentelemetry-collector-builder

COPY otel-config.yaml /otel-config.yaml

CMD /tmp/dist/otelcol-custom --config=/otel-config.yaml 

