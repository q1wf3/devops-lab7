FROM --platform=$BUILDPLATFORM golang:1.22.12 AS builder
WORKDIR /src
COPY ticket-service/go.mod ticket-service/go.sum ./
RUN go mod download
COPY ticket-service/ ./
ARG TARGETOS
ARG TARGETARCH
RUN CGO_ENABLED=0 GOOS=${TARGETOS:-linux} GOARCH=${TARGETARCH} \
    go build -trimpath -ldflags="-s -w" -o /out/server ./order

FROM alpine:3.20
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=builder /out/server ./server
USER 65532:65532
EXPOSE 8080 9090
ENTRYPOINT ["/app/server"]
