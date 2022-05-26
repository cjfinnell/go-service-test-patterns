ARG GO_VERSION=1.18
ARG ALPINE_VERSION=3.16

FROM golang:${GO_VERSION}-alpine${ALPINE_VERSION} as builder
WORKDIR /build
ADD . .
RUN go build -o app

FROM golang:${GO_VERSION} as tester
WORKDIR /test
ADD . .
CMD make test.integration

FROM alpine:${ALPINE_VERSION} as release
WORKDIR /app
COPY --from=builder /build/app .
RUN adduser -S -D -H -h /app appuser
USER appuser
EXPOSE 8080
CMD ./app
