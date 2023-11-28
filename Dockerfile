FROM golang:1.20 as builder

WORKDIR /app
COPY . /app

ENV CGO_ENABLED=0
ENV GO111MODULE=on
RUN make tools
RUN make build

FROM alpine:3

RUN apk add --no-cache ca-certificates
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

RUN chown appuser:appuser /server
RUN chown appuser:appuser /server/*
USER appuser:appuser

ENTRYPOINT ["/app/approver"]