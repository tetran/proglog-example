FROM --platform=linux/amd64 golang:1.21-alpine AS build
WORKDIR /go/src/proglog-example
COPY . .
RUN CGO_ENABLED=0 go build -o /go/bin/proglog-example ./cmd/proglog-example

RUN GRPC_HEALTH_PROBE_VERSION=v0.4.8 && \
    wget -qO/go/bin/grpc_health_probe \
    https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 && \
    chmod +x /go/bin/grpc_health_probe

FROM --platform=linux/amd64 alpine
COPY --from=build /go/bin/proglog-example /bin/proglog-example
COPY --from=build /go/bin/grpc_health_probe /bin/grpc_health_probe
ENTRYPOINT [ "/bin/proglog-example" ]
