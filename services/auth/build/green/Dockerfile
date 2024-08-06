FROM golang:1.22.5-alpine AS builder

WORKDIR /go/src/app

COPY auth/cmd/app/green/main.go .

RUN go build -o /go/bin/app main.go

FROM alpine:3.20.2

COPY --from=builder /go/bin/app /usr/local/bin/app

CMD ["app"]

EXPOSE 8080
