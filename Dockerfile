FROM golang:1.23 AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY cmd ./cmd
COPY internal ./internal
RUN GOOS=linux GOARCH=amd64 go build --ldflags '-extldflags "-static"' -o ./deepwildcard.bin ./cmd/deepwildcard 

FROM scratch
COPY --from=builder /app/deepwildcard.bin /usr/bin/deepwildcard
COPY config.yaml /etc/deepwildcard/config.yaml
EXPOSE 9000
VOLUME [ "/etc/deepwildcard" ]
ENTRYPOINT [ "/usr/bin/deepwildcard" ]