FROM golang:1.20 AS builder
ARG SERVICE
WORKDIR /go/src/github.com/alexanderjophus/kie
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -o app ./svc/${SERVICE}

FROM gcr.io/distroless/base
COPY --from=builder /go/src/github.com/alexanderjophus/kie/app ./app
CMD ["./app"]