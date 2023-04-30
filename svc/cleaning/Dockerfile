FROM golang:1.20 AS builder
WORKDIR /go/src/github.com/alexanderjophus/kie/svc/cleaning
COPY . /go/src/github.com/alexanderjophus/kie/svc/cleaning
RUN CGO_ENABLED=0 GOOS=linux go build -a -ldflags "-s -w -X github.com/alexanderjophus/kie/svc/cleaning/cmd.inDir=/pfs/statsapi/api/v1/people/ -X github.com/alexanderjophus/kie/svc/cleaning/cmd.outDir=/pfs/out/" -o app .

FROM gcr.io/distroless/base
COPY --from=builder /go/src/github.com/alexanderjophus/kie/svc/cleaning/app ./app
CMD ["./app"]