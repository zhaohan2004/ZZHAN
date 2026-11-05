FROM golang:1.26-alpine AS build
WORKDIR /src
COPY ZZHAN/go.mod ZZHAN/go.sum ./
RUN go mod download
COPY ZZHAN/cmd ./cmd
COPY ZZHAN/internal ./internal
COPY ZZHAN/pkg ./pkg
RUN go test ./...
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/server ./cmd/server

FROM alpine:3.23
RUN apk add --no-cache ca-certificates tzdata \
    && addgroup -S app && adduser -S -G app app \
    && mkdir -p /app/logs /app/pkg/uploads \
    && chown -R app:app /app
WORKDIR /app
COPY --from=build /out/server /app/server
USER app
EXPOSE 9900
HEALTHCHECK --interval=15s --timeout=5s --start-period=30s --retries=4 \
  CMD wget -q -O /dev/null http://127.0.0.1:9900/api/v1/health || exit 1
ENTRYPOINT ["/app/server"]
CMD ["-c", "/app/config.yaml"]
