FROM golang:1.16-alpine AS builder

WORKDIR /usr/src/app

COPY . .

RUN go build -o app .

FROM alpine:latest AS runner

WORKDIR /usr/src/app

COPY --from=builder /usr/src/app/app ./app

CMD ["./app"]