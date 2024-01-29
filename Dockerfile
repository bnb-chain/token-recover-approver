FROM golang:1.20-alpine as builder

WORKDIR /app
COPY . /app

RUN apk add --no-cache make git bash build-base linux-headers libc-dev
ENV GO111MODULE=on
ENV CGO_ENABLED=1
ENV CGO_CFLAGS="-D_LARGEFILE64_SOURCE"
RUN make tools
RUN make build

FROM alpine:3

RUN apk add --no-cache ca-certificates libstdc++
RUN mkdir -p /app
WORKDIR /app

COPY --from=builder /app/build/bin/approver /app/approver

# Create appuser.
ENV USER=appuser
ENV UID=1001

RUN adduser \
--disabled-password \
--gecos "application user" \
--no-create-home \
--uid "${UID}" \
"${USER}"

RUN chown appuser:appuser /app
RUN chown appuser:appuser /app/*
USER appuser:appuser

ENTRYPOINT ["/app/approver"]