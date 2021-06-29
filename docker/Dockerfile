#FROM golang:1.16.5-buster
FROM golang:1.15-buster
 
#RUN cd ~ && git clone https://github.com/open-telemetry/opentelemetry-collector-builder.git && pwd && cd opentelemetry-collector-builder && pwd && ls && git checkout tags/v0.29.0 && echo "here" && ls

#RUN GO111MODULE=on go get github.com/open-telemetry/opentelemetry-collector-builder
#RUN GO111MODULE=on go get github.com/open-telemetry/opentelemetry-collector-builder@v0.25.0
RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder
#RUN GO111MODULE=on go get github.com/observatorium/opentelemetry-collector-builder@v0.1.2

RUN git clone https://github.com/librato/opentelemetry-collector-contrib.git && cd opentelemetry-collector-contrib && git checkout AO-18760-ao-exporter && cd ..

RUN git clone https://github.com/appoptics/appoptics-apm-go.git && cd appoptics-apm-go && git checkout AO-18760-ao-exporter && cd ..

COPY .otelcol-builder.yaml /root/.otelcol-builder.yaml

RUN opentelemetry-collector-builder

COPY otel-config.yaml /otel-config.yaml

CMD /tmp/dist/otelcol-custom --config=/otel-config.yaml 

