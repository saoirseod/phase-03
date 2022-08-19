FROM amd64/golang:alpine AS builder

WORKDIR /hello/bin
# RUN GRPC_HEALTH_PROBE_VERSION=v0.3.6 && \
#     wget -qO ./grpc_health_probe https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64
# RUN chmod +x ./grpc_health_probe
WORKDIR /hello

ADD . ./

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o bin/server server/server.go

FROM scratch
# RUN   apt update \
#   &&   apt add ca-certificates \
#   &&   update-ca-certificates
WORKDIR /
COPY --from=builder /hello/bin/server .
# COPY --from=builder /hello/bin/grpc_health_probe .
EXPOSE 9000
ENTRYPOINT ["/server"]

