FROM golang:1.21.6-alpine AS builder

WORKDIR /app

COPY ./ ./

RUN go version
ENV GOPATH=/

RUN apk update && apk add --no-cache postgresql-client

RUN go mod download
RUN go build -o avito_testcase ./cmd/main.go

FROM alpine:latest

WORKDIR /root/

COPY --from=builder /app .

RUN ls -la

RUN chmod +x avito_testcase

CMD ["./avito_testcase"]