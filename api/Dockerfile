FROM golang:1.24-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -trimpath -ldflags="-s -w" -o /out/arumonogohan-api .

FROM alpine:3.22

RUN adduser -D -H -u 10001 appuser
COPY --from=builder /out/arumonogohan-api /arumonogohan-api

USER appuser
EXPOSE 8080
ENTRYPOINT ["/arumonogohan-api"]
