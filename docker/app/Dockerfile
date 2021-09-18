FROM golang:latest as builder
WORKDIR /workdir
COPY src .
RUN go mod tidy
# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Build
RUN go build -o app

# Runtime Container
FROM alpine
RUN apk add --no-cache ca-certificates
COPY --from=builder /workdir/app /app
COPY --from=builder /workdir/public /public
RUN mkdir -p /tmp/images

ENTRYPOINT ["/app"]
