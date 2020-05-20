FROM golang:1.14 AS builder
WORKDIR /go/src/github.com/polyse/frontend-service

COPY go.mod .
COPY go.sum .
RUN go mod download

COPY . .
RUN GOOS=linux CGO_ENABLED=0 go build -installsuffix cgo -o app github.com/polyse/frontend-service/cmd/front

FROM alpine:latest
RUN apk --no-cache add ca-certificates
WORKDIR /app
COPY --from=0 /go/src/github.com/polyse/frontend-service .

ENV LOG_FMT json
ENV LISTEN 0.0.0.0:9900

EXPOSE 9900

ENTRYPOINT ["/app/app"]
