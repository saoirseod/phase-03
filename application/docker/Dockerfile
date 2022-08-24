FROM amd64/golang:alpine AS builder

WORKDIR /hello/bin

WORKDIR /hello

ADD . ./

RUN go mod download

RUN GOOS=linux GOARCH=amd64 go build -o bin/server server/server.go

FROM scratch

WORKDIR /
COPY --from=builder /hello/bin/server .

EXPOSE 9000
ENTRYPOINT ["/server"]

