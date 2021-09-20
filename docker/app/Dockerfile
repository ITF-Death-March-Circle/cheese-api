FROM notchman/opencv:latest as builder-cv
WORKDIR /cheese

COPY /src/face_detect/ ./
# COPY CMakeLists.txt .

RUN cmake . && make

FROM golang:latest as builder
WORKDIR /workdir
COPY src ./
RUN go mod tidy
# Set Environment Variable
ENV CGO_ENABLED=0
ENV GOOS=linux
ENV GOARCH=amd64
# Build
RUN go build -o app

FROM golang:bullseye
COPY --from=builder /workdir/app /app
COPY --from=builder-cv /cheese/DisplayImage /DisplayImage
COPY --from=builder-cv /usr/local/lib /usr/local/lib
# テスト
# RUN mkdir /cheese/images -p
# RUN ./DisplayImage
ENTRYPOINT ["/app"]
