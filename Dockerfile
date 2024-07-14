ARG GOLANG_VERSION=1.22-alpine
FROM golang:${GOLANG_VERSION} AS base
WORKDIR /app

FROM base AS builder

COPY go.mod ./
RUN go mod download

COPY *.go ./
RUN CGO_ENABLED=0 go build -o echo

FROM gcr.io/distroless/static-debian12
COPY --from=builder /app/echo /app/echo
ENTRYPOINT [ "/app/echo" ]
