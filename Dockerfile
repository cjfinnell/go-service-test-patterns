FROM golang:1.18-alpine3.16 as builder
WORKDIR /build
ADD . .
RUN go build -o app

FROM golang:1.18 as tester
WORKDIR /test
ADD . .
CMD make test.integration

FROM alpine:3.16 as release
WORKDIR /app
COPY --from=builder /build/app .
RUN adduser -S -D -H -h /app appuser
USER appuser
EXPOSE 8080
CMD ./app
